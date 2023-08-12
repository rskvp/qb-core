package qb_log

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

const (
	outConsole = "console"
	outFile    = "file"
)

type logMessage struct {
	level Level
	args  []interface{}
}

type Logger struct {
	filename       string
	level          Level
	mux            sync.Mutex
	channel        chan *logMessage
	initialized    bool
	datePattern    string
	messagePattern string
	rotateEnable   bool
	rotateSizeMb   float64
	rotateDir      string
	rotateMaxFiles int
	outputs        []string
}

func NewLogger() *Logger {
	instance := new(Logger)
	instance.datePattern = DateFormatStandard
	instance.messagePattern = pattern
	instance.rotateMaxFiles = 10
	instance.RotateMaxSizeMb(1)
	instance.outputs = make([]string, 0)

	// init the buffered channel
	instance.channel = make(chan *logMessage, 10)
	go instance.receive(instance.channel)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Logger) String() string {
	m := map[string]interface{}{
		"level":        instance.level.String(),
		"outputs":      instance.outputs,
		"filename":     instance.filename,
		"rotate":       instance.rotateEnable,
		"rotate_limit": instance.rotateSizeMb,
	}
	return qb_utils.JSON.Stringify(m)
}

func (instance *Logger) GetLevel() Level {
	if nil != instance {
		return instance.level
	}
	return InfoLevel
}

func (instance *Logger) SetLevel(v interface{}) *Logger {
	if nil != instance {
		if level, ok := v.(Level); ok {
			instance.level = level
		} else if level, ok := v.(*Level); ok {
			instance.level = *level
		} else if level, ok := v.(string); ok {
			instance.level, _ = ParseLevel(level)
		}
	}
	return instance
}

func (instance *Logger) GetFilename() string {
	if nil != instance {
		return instance.filename
	}
	return ""
}

func (instance *Logger) SetFilename(filename string) *Logger {
	if nil != instance {
		if len(filename) == 0 {
			filename = "logging.log"
		}
		root := qb_utils.Paths.WorkspacePath("./logging")
		instance.filename = qb_utils.Paths.Absolutize(filename, root)

		instance.rotateDir = qb_utils.Paths.Concat(qb_utils.Paths.Dir(instance.filename), "log-rotate")
	}
	return instance
}

func (instance *Logger) OutConsole(value bool) *Logger {
	if nil != instance {
		if value {
			instance.outputs = qb_utils.Arrays.AppendUnique(instance.outputs, outConsole).([]string)
		} else {
			instance.outputs = qb_utils.Arrays.Remove(outConsole, instance.outputs).([]string)
		}
	}
	return instance
}

func (instance *Logger) OutFile(value bool) *Logger {
	if nil != instance {
		if value {
			instance.outputs = qb_utils.Arrays.AppendUnique(instance.outputs, outFile).([]string)
		} else {
			instance.outputs = qb_utils.Arrays.Remove(outFile, instance.outputs).([]string)
		}
	}
	return instance
}

func (instance *Logger) RotateEnable(value bool) *Logger {
	if nil != instance {
		instance.rotateEnable = value
	}
	return instance
}

func (instance *Logger) RotateMaxSizeMb(value float64) *Logger {
	if nil != instance {
		instance.rotateSizeMb = value
		instance.rotateEnable = value > 0
	}
	return instance
}

func (instance *Logger) SetDateFormat(format string) *Logger {
	if nil != instance {
		instance.datePattern = format
	}
	return instance
}

func (instance *Logger) GetDateFormat() (response string) {
	if nil != instance {
		response = instance.datePattern
	}
	return
}

func (instance *Logger) SetMessageFormat(format string) *Logger {
	if nil != instance {
		instance.messagePattern = format
	}
	return instance
}

func (instance *Logger) GetMessageFormat() (response string) {
	if nil != instance {
		response = instance.messagePattern
	}
	return
}

func (instance *Logger) Close() {
	if nil != instance.channel {
		close(instance.channel)
	}
}

func (instance *Logger) Flush() {
	if nil != instance {
		// TODO: add a flush function
		time.Sleep(1 * time.Second)
	}
}

func (instance *Logger) Panic(args ...interface{}) {
	message := new(logMessage)
	message.level = PanicLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Panicf(message string, args ...interface{}) {
	if nil != instance {
		instance.Panic(fmt.Sprintf(message, args...))
	}
}

func (instance *Logger) Error(args ...interface{}) {
	message := new(logMessage)
	message.level = ErrorLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Errorf(message string, args ...interface{}) {
	if nil != instance {
		instance.Error(fmt.Sprintf(message, args...))
	}
}

func (instance *Logger) Warn(args ...interface{}) {
	message := new(logMessage)
	message.level = WarnLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Warnf(message string, args ...interface{}) {
	if nil != instance {
		instance.Warn(fmt.Sprintf(message, args...))
	}
}

func (instance *Logger) Info(args ...interface{}) {
	message := new(logMessage)
	message.level = InfoLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Infof(message string, args ...interface{}) {
	if nil != instance {
		instance.Info(fmt.Sprintf(message, args...))
	}
}

func (instance *Logger) Debug(args ...interface{}) {
	message := new(logMessage)
	message.level = DebugLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Debugf(message string, args ...interface{}) {
	if nil != instance {
		instance.Debug(fmt.Sprintf(message, args...))
	}
}

func (instance *Logger) Trace(args ...interface{}) {
	message := new(logMessage)
	message.level = TraceLevel
	message.args = args
	if nil != instance.channel {
		instance.channel <- message
	}
}

func (instance *Logger) Tracef(message string, args ...interface{}) {
	if nil != instance {
		instance.Trace(fmt.Sprintf(message, args...))
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Logger) initialize() {
	if !instance.initialized {
		instance.initialized = true

		if len(instance.outputs) == 0 {
			instance.outputs = append(instance.outputs, outFile) // default is on file
		}

		if instance.canLogToFile() {
			// set logging file
			if len(instance.filename) == 0 {
				instance.SetFilename("") // default filename
			}

			_ = qb_utils.Paths.Mkdir(instance.filename)

			if ok, _ := qb_utils.Paths.Exists(instance.filename); ok {
				instance.rotate(instance.filename)
			}
		}
	}
}

func (instance *Logger) receive(ch <-chan *logMessage) {
	// loop until channel is open
	for message := range ch {
		instance.doLog(message.level, message.args...)
	}
}

func (instance *Logger) doLog(level Level, args ...interface{}) {
	instance.initialize()

	if b := instance.level < level; b {
		return
	}
	var buf strings.Builder
	for i, arg := range args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", arg))
	}
	buf.WriteString("\n")

	m := map[string]interface{}{
		"level":   strings.ToUpper(level.String()),
		"message": buf.String(),
	}
	if len(instance.datePattern) > 0 {
		m["date"] = qb_utils.Dates.FormatDate(time.Now(), instance.datePattern)
	}
	message, err := qb_utils.Formatter.MergeText(instance.messagePattern, m)
	if nil != err {
		log.Panicln(fmt.Sprintf("Panic error formatting log message: '%s'", err))
	}

	instance.writeOutput(message)
}

func (instance *Logger) writeOutput(text string) {
	// PANIC RECOVERY
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			message := qb_utils.Strings.Format("[panic] logger.writeOutput('%s'): '%s'", text, r)
			log.Println(message)
		}
	}()

	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		if instance.canLogToFile() {
			instance.rotateVerify(instance.filename)

			err := writeToFile(text, instance.filename)
			if nil != err {
				log.Panicln(fmt.Sprintf("Panic error writing log file '%s': '%s'", instance.filename, err))
			}
		}

		if instance.canLogToConsole() {
			fmt.Print(text)
		}
	}
}

func (instance *Logger) canLogToFile() bool {
	if nil != instance {
		return len(instance.filename) > 0 && qb_utils.Arrays.IndexOf(outFile, instance.outputs) > -1
	}
	return false
}

func (instance *Logger) canLogToConsole() bool {
	if nil != instance {
		return len(instance.filename) == 0 || qb_utils.Arrays.IndexOf(outConsole, instance.outputs) > -1
	}
	return false
}

func (instance *Logger) rotateVerify(filename string) {
	if instance.rotateEnable && instance.rotateSizeMb > 0 {
		if ok, _ := qb_utils.Paths.Exists(filename); ok {
			size, _ := qb_utils.IO.FileSize(filename)
			if size > 0 && qb_utils.Convert.ToMegaBytes(size) > instance.rotateSizeMb {
				instance.rotate(filename)
			}
		}
	}
}

func (instance *Logger) rotate(filename string) {
	if instance.rotateEnable {
		if exists, _ := qb_utils.Paths.Exists(instance.rotateDir + qb_utils.OS_PATH_SEPARATOR); !exists {
			_ = qb_utils.Paths.Mkdir(instance.rotateDir + qb_utils.OS_PATH_SEPARATOR)
		}
		newFilename := qb_utils.Paths.ChangeFileNameWithSuffix(filename, qb_utils.Dates.FormatDate(time.Now(), "-yyyy-MM-dd-HH-mm-ss"))
		newFilename = qb_utils.Paths.Concat(instance.rotateDir, qb_utils.Paths.FileName(newFilename, true))
		_, _ = qb_utils.IO.CopyFile(filename, newFilename)
		_, _ = qb_utils.IO.WriteTextToFile("", instance.filename)
		// list files
		if instance.rotateMaxFiles > 0 {
			files, _ := qb_utils.Paths.ListFiles(instance.rotateDir, "*.log")
			if len(files) > instance.rotateMaxFiles {
				_ = qb_utils.IO.Remove(files[0])
			}
		}
	}
}

func writeToFile(text, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm) // 0644
	if err != nil {
		return err
	}

	defer file.Close()
	w := bufio.NewWriter(file)
	_, err = w.WriteString(text)
	_ = w.Flush()
	return err
}
