package qb_sys

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rskvp/qb-core/qb_"
)

type ShutdownCallback func(ctx context.Context) error

//----------------------------------------------------------------------------------------------------------------------
//	SysHelper
//----------------------------------------------------------------------------------------------------------------------

type SysHelper struct {
}

var Sys *SysHelper

func init() {
	Sys = new(SysHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *SysHelper) GetInfo() *InfoObject {
	info := getInfo()

	// memory
	info.MemoryUsage = NewMemoryUsageInfo().ToString()

	return info
}

func (instance *SysHelper) GetOS() string {
	return runtime.GOOS
}

func (instance *SysHelper) IsMac() bool {
	return "darwin" == instance.GetOS()
}

func (instance *SysHelper) IsLinux() bool {
	return "linux" == instance.GetOS()
}

func (instance *SysHelper) IsWindows() bool {
	return "windows" == instance.GetOS()
}

func (instance *SysHelper) GetOSVersion() string {
	return instance.GetInfo().Core
}

// Shutdown  the machine
func (instance *SysHelper) Shutdown(a ...string) error {
	adminPsw := ""
	if len(a) == 1 {
		adminPsw = a[0]
	}
	return shutdown(adminPsw)
}

func (instance *SysHelper) OnSignal(callback func(s os.Signal), signals ...os.Signal) (err error) {
	// signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	channel := make(chan os.Signal, 1)
	if len(signals) == 0 {
		err = errors.New("missing_signals")
		return
	}
	if nil == callback {
		err = errors.New("missing_callback")
		return
	}
	signal.Notify(channel, signals...)
	go func() {
		// RECOVER
		if r := recover(); r != nil {
			// do nothing, but avoid propagation
		}
		// wait
		s := <-channel
		// handler
		callback(s)
	}()
	return
}

func (instance *SysHelper) GracefulShutdownBackground(actions map[string]ShutdownCallback, logger qb_.ILogger) chan struct{} {
	return instance.GracefulShutdownWithContext(context.Background(), time.Minute*1, actions, logger)
}

func (instance *SysHelper) GracefulShutdownWithContext(ctx context.Context, timeout time.Duration, actions map[string]ShutdownCallback, logger qb_.ILogger) chan struct{} {
	// creates a wait channel
	wait := make(chan struct{})
	_ = instance.OnSignal(func(s os.Signal) {
		msg := "STARTING SHUTDOWN..."
		if nil != logger {
			logger.Info(msg)
		} else {
			log.Println(msg)
		}

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			msg := fmt.Sprintf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			if nil != logger {
				logger.Info(msg)
			} else {
				log.Println(msg)
			}
			os.Exit(0)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		// Do the operations asynchronously to save time
		for key, op := range actions {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				msg := fmt.Sprintf("\tcleaning up: %s", innerKey)
				if nil != logger {
					logger.Info(msg)
				} else {
					log.Println(msg)
				}
				if err := innerOp(ctx); err != nil {
					msg = fmt.Sprintf("\t%s: clean up failed: %s", innerKey, err.Error())
					if nil != logger {
						logger.Info(msg)
					} else {
						log.Println(msg)
					}
					return
				}

				msg = fmt.Sprintf("\t%s was shutdown gracefully", innerKey)
				if nil != logger {
					logger.Info(msg)
				} else {
					log.Println(msg)
				}
			}()
		}
		wg.Wait()

		// close wait channel
		close(wait)
		msg = "TERMINATED SHUTDOWN."
		if nil != logger {
			logger.Info(msg)
		} else {
			log.Println(msg)
		}
	}, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGQUIT)
	return wait // return wait channel
}

// ID returns the platform specific machine id of the current host OS.
// Regard the returned id as "confidential" and consider using ProtectedID() instead.
// THANKS TO: github.com/denisbrodbeck/machineid
func (instance *SysHelper) ID() (string, error) {
	id, err := machineID()
	if err != nil {
		return "", fmt.Errorf("machineid: %v", err)
	}
	return id, nil
}

// ProtectedID returns a hashed version of the machine ID in a cryptographically secure way,
// using a fixed, application-specific key.
// Internally, this function calculates HMAC-SHA256 of the application ID, keyed by the machine ID.
// THANKS TO: github.com/denisbrodbeck/machineid
func (instance *SysHelper) ProtectedID(appID string) (string, error) {
	id, err := instance.ID()
	if err != nil {
		return "", fmt.Errorf("machineid: %v", err)
	}
	return _protect(appID, id), nil
}

func (instance *SysHelper) FindCurrentProcess() (*os.Process, error) {
	return os.FindProcess(syscall.Getpid())
}

func (instance *SysHelper) KillCurrentProcess() error {
	p, err := instance.FindCurrentProcess()
	if nil != err {
		return err
	}
	return p.Signal(os.Interrupt)
}

func (instance *SysHelper) KillProcessByPid(pid int) error {
	p, err := os.FindProcess(pid)
	if nil != err {
		return err
	}
	err = p.Kill()
	if nil != err {
		// process found but already closed?
	}
	return err
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type InfoObject struct {
	GoOS        string `json:"goos"`
	Kernel      string `json:"kernel"`
	Core        string `json:"core"`
	Platform    string `json:"platform"`
	OS          string `json:"os"`
	Hostname    string `json:"hostname"`
	CPUs        int    `json:"cpus"`
	MemoryUsage string `json:"memory_usage"`
}

func (instance *InfoObject) VarDump() {
	fmt.Println("GoOS:", instance.GoOS)
	fmt.Println("Kernel:", instance.Kernel)
	fmt.Println("Core:", instance.Core)
	fmt.Println("Platform:", instance.Platform)
	fmt.Println("OS:", instance.OS)
	fmt.Println("Hostname:", instance.Hostname)
	fmt.Println("CPUs:", instance.CPUs)
	fmt.Println("MemoryUsage:", instance.MemoryUsage)
}

func (instance *InfoObject) ToString() string {
	return fmt.Sprintf("GoOS:%v,Kernel:%v,Core:%v,Platform:%v,OS:%v,Hostname:%v,CPUs:%v, MemoryUsage:%v", instance.GoOS, instance.Kernel, instance.Core, instance.Platform, instance.OS, instance.Hostname, instance.CPUs, instance.MemoryUsage)
}

func (instance *InfoObject) ToJsonString() string {
	return _stringify(instance)
}

type MemObject struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
	HeapSys    uint64 `json:"heap_sys"`
}

func NewMemoryUsageInfo() *MemObject {
	instance := new(MemObject)
	// memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	instance.Alloc = m.Alloc
	instance.TotalAlloc = m.TotalAlloc
	instance.Sys = m.Sys
	instance.NumGC = m.NumGC
	instance.HeapSys = m.HeapSys

	return instance
}

func (instance *MemObject) String() string {
	return _stringify(instance)
}

func (instance *MemObject) ToString() string {
	return fmt.Sprintf("Alloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v, HeapSys = %v",
		_formatBytes(instance.Alloc),
		_formatBytes(instance.TotalAlloc),
		_formatBytes(instance.Sys), instance.NumGC,
		_formatBytes(instance.HeapSys))

}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// _run wraps `exec.Command` with easy access to stdout and stderr.
func _run(stdout, stderr io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}

// _protect calculates HMAC-SHA256 of the application ID, keyed by the machine ID and returns a hex-encoded string.
func _protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return hex.EncodeToString(mac.Sum(nil))
}

func _readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func _trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

func _bytes(entity interface{}) []byte {
	if nil != entity {
		if s, b := entity.(string); b {
			return []byte(s)
		}
		b, err := json.Marshal(&entity)
		if nil == err {
			return b
		}
	}

	return []byte{}
}

func _stringify(entity interface{}) string {
	if s, b := entity.(string); b {
		if strings.Index(s, "\"") != 0 {
			// quote to not quoted string
			return strconv.Quote(_toString(s))
		}
	}
	return string(_bytes(entity))
}

func _toString(val interface{}) string {
	if nil == val {
		return ""
	}
	// string
	s, ss := val.(string)
	if ss {
		return s
	}
	// integer
	i, ii := val.(int)
	if ii {
		return strconv.Itoa(i)
	}
	// float32
	f, ff := val.(float32)
	if ff {
		return fmt.Sprintf("%g", f) // Exponent as needed, necessary digits only
	}
	// float 64
	F, FF := val.(float64)
	if FF {
		return fmt.Sprintf("%g", F) // Exponent as needed, necessary digits only
		// return strconv.FormatFloat(F, 'E', -1, 64)
	}

	// boolean
	b, bb := val.(bool)
	if bb {
		return strconv.FormatBool(b)
	}

	// array
	if aa, _ := _isArray(val); aa {
		// byte array??
		if ba, b := val.([]byte); b {
			return string(ba)
		} else {
			data, err := json.Marshal(val)
			if nil == err {
				return string(data)
			}
		}
	}

	// map
	if b, _ := _isMap(val); b {
		data, err := json.Marshal(val)
		if nil == err {
			return string(data)
		}
	}

	// struct
	if b, _ := _isStruct(val); b {
		data, err := json.Marshal(val)
		if nil == err {
			return string(data)
		}
	}

	// undefined value
	return fmt.Sprintf("%v", val)
}

func _toInt64Def(val interface{}, defVal int64) int64 {
	switch i := val.(type) {
	case float32:
		return int64(i)
	case float64:
		return int64(i)
	case int:
		return int64(i)
	case int8:
		return int64(i)
	case int16:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return i
	case uint8:
		return int64(i)
	case uint16:
		return int64(i)
	case uint32:
		return int64(i)
	case uint64:
		return int64(i)
	case string:
		v, err := strconv.ParseInt(i, 10, 64)
		if nil == err {
			return v
		}
	}
	return defVal
}

func _isStruct(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Struct:
		return true, rt
	case reflect.Ptr:
		return _isStruct(rt.Elem().Interface())
	default:
		return false, rt
	}
}

func _isMap(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Map:
		return true, rt
	case reflect.Ptr:
		return _isMap(rt.Elem().Interface())
	default:
		return false, rt
	}
}

func _isArray(val interface{}) (bool, reflect.Value) {
	rt := reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return true, rt
	default:
		return false, rt
	}
}

func _formatBytes(i interface{}) string {
	n := uint64(_toInt64Def(i, 0))
	return _fmtBytes(n)
}

func _fmtBytes(b uint64) string {
	const unit = 1024
	if b < uint64(1024) {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
