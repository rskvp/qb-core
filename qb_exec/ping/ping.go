package ping

import (
	"io"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_sys"
	"github.com/rskvp/qb-core/qb_utils"
)

var pingCommand = "ping"
var pingOnceParams = "-c 1"

type PingResponse struct {
	Header     string                  `json:"header"`
	Body       []string                `json:"body"`
	Footer     string                  `json:"footer"`
	Rows       []*PingResponseRow      `json:"rows"`
	Statistics *PingResponseStatistics `json:"statistics"`
}

func (instance *PingResponse) String() string {
	return qb_utils.JSON.Stringify(instance)
}

type PingResponseRow struct {
	Bytes  int     `json:"bytes"`
	TimeMs float32 `json:"time_ms"`
	Ttl    string  `json:"ttl"`
	Source string  `json:"source"`
}

func (instance *PingResponseRow) String() string {
	return qb_utils.JSON.Stringify(instance)
}

type PingResponseStatistics struct {
	PacketTransmitted int     `json:"packet_transmitted"`
	PacketReceived    int     `json:"packet_received"`
	PacketLoss        float32 `json:"packet_loss"`
	Min               float32 `json:"min"`
	Max               float32 `json:"max"`
	Avg               float32 `json:"avg"`
}

func (instance *PingResponseStatistics) String() string {
	return qb_utils.JSON.Stringify(instance)
}

type PingExec struct {
	dirWork string
}

func NewPingExec() *PingExec {
	instance := new(PingExec)
	instance.dirWork = qb_utils.Paths.Absolute("./")
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *PingExec) SetDir(dir string) {
	instance.dirWork = dir
}

func (instance *PingExec) GetDir() string {
	return instance.dirWork
}

func (instance *PingExec) PingOut(target string, out io.Writer) (*executor.ConsoleProgramSession, error) {
	if nil != instance {
		program := instance.program()
		program.OutWriterAppend(out)
		program.ErrorWriterAppend(out)
		return program.RunAsync(target)
	}
	return nil, nil
}

func (instance *PingExec) PingOnce(target string) (*PingResponse, error) {
	if nil != instance {
		session, err := instance.program().Run(pingOnceParams, target)
		if nil != err {
			return nil, err
		}
		if nil != session {
			out := session.StdOut()
			pr := instance.parseResponse(out)
			pr.Statistics = instance.summarize(pr.Rows)
			return pr, nil
		}
	}
	return nil, nil
}

func (instance *PingExec) PingCount(target string, n int) (*PingResponse, error) {
	response := new(PingResponse)
	if nil != instance {
		for i := 0; i < n; i++ {
			session, err := instance.program().Run(pingOnceParams, target)
			if nil != err {
				return nil, err
			}
			if nil != session {
				out := session.StdOut()
				pr := instance.parseResponse(out)
				if len(response.Header) == 0 {
					response.Header = pr.Header
					response.Footer = pr.Footer
				}
				response.Body = append(response.Body, pr.Body...)
				response.Rows = append(response.Rows, pr.Rows...)
			}
			time.Sleep(1 * time.Second)
		}
	}
	response.Statistics = instance.summarize(response.Rows)
	return response, nil
}

func (instance *PingExec) Ping(target string, exitCallback func(row *PingResponseRow) bool) (*PingResponse, error) {
	response := new(PingResponse)
	if nil != instance && nil != exitCallback {
		for {
			session, err := instance.program().Run(pingOnceParams, target)
			if nil != err {
				return nil, err
			}
			if nil != session {
				out := session.StdOut()
				pr := instance.parseResponse(out)
				if len(response.Header) == 0 {
					response.Header = pr.Header
					response.Footer = pr.Footer
				}
				response.Body = append(response.Body, pr.Body...)
				response.Rows = append(response.Rows, pr.Rows...)
				if exitCallback(pr.Rows[0]) {
					break
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
	response.Statistics = instance.summarize(response.Rows)
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *PingExec) program() *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(pingCommand, instance.dirWork)
}

func (instance *PingExec) parseResponse(raw string) *PingResponse {
	response := new(PingResponse)
	tokens := strings.Split(raw, "\n\n")
	if len(tokens) == 2 {
		response.Footer = tokens[1]
		rows := strings.Split(tokens[0], "\n")
		if len(rows) > 1 {
			response.Header = rows[0]
			response.Body = rows[1:]
			for _, row := range response.Body {
				response.Rows = append(response.Rows, ParsePingRow(row))
			}
		}
	}
	return response
}

func (instance *PingExec) summarize(rows []*PingResponseRow) *PingResponseStatistics {
	statistics := new(PingResponseStatistics)
	count := float32(len(rows))
	total := float32(0.0)
	for _, row := range rows {
		total += row.TimeMs
		if statistics.Min == 0 || statistics.Min > row.TimeMs {
			statistics.Min = row.TimeMs
		}
		if statistics.Max < row.TimeMs {
			statistics.Max = row.TimeMs
		}
		if row.TimeMs > 0 {
			statistics.PacketReceived += 1
		}
	}
	statistics.PacketTransmitted = len(rows)
	statistics.PacketLoss = 100 - float32(statistics.PacketTransmitted/statistics.PacketReceived)*100
	statistics.Avg = total / count
	return statistics
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func init() {
	pingCommand = "ping"
	if qb_sys.Sys.IsWindows() {
		pingOnceParams = "-n 1"
	}
}

func ParsePingRow(row string) *PingResponseRow {
	// WIN > 	Risposta da 216.58.206.36: byte=32 durata=11ms TTL=128
	// MAC > 	64 bytes from 216.58.198.36: icmp_seq=0 ttl=119 time=11.090 ms
	// LINUX > 	64 bytes from arn11s11-in-x04.1e100.net (2a00:1450:400f:804::2004): icmp_seq=1 ttl=117 time=27.3 ms
	response := new(PingResponseRow)
	response.Source = row
	response.Bytes = 64
	tokens := qb_utils.Strings.SplitLast(row, ':')
	if len(tokens) == 2 {
		if qb_sys.Sys.IsWindows() {
			// windows "Risposta da 216.58.206.36: byte=32 durata=11ms TTL=128"
			data := qb_utils.Strings.SplitTrimSpace(tokens[1], " ")
			for _, s := range data {
				if len(s) > 0 {
					kp := qb_utils.Strings.SplitTrimSpace(s, "=")
					if len(kp) == 2 {
						switch kp[0] {
						case "byte", "bytes":
							response.Bytes = qb_utils.Convert.ToInt(kp[1])
						case "TTL", "ttl":
							response.Ttl = kp[1]
						default:
							response.TimeMs = qb_utils.Convert.ToFloat32(strings.ReplaceAll(kp[1], "ms", ""))
						}
					}
				}
			}
		} else {
			// other os
			data := qb_utils.Strings.SplitTrimSpace(tokens[1], " ")
			for _, s := range data {
				if len(s) > 0 {
					kp := qb_utils.Strings.SplitTrimSpace(s, "=")
					if len(kp) == 2 {
						switch kp[0] {
						case "byte", "bytes":
							response.Bytes = qb_utils.Convert.ToInt(kp[1])
						case "TTL", "ttl":
							response.Ttl = kp[1]
						case "icmp_seq":
							response.Ttl = kp[1]
						default:
							response.TimeMs = qb_utils.Convert.ToFloat32(kp[1])
						}
					}
				}
			}
		}
	}
	return response
}
