package qbc

import (
	"errors"
	"fmt"
	"log"

	"github.com/rskvp/qb-core/qb_log"
)

const Version = "0.0.1"

const (
	ModeProduction = "production"
	ModeDebug      = "debug"
)

var (
	AppName    = "New Application"
	AppVersion = "0.0.1" // change this for a global AppName variable

	ErrorSystem = errors.New("panic_system_error")
)

func Recover(args ...interface{}) {
	if r := recover(); r != nil {
		// recovered from panic
		message := getMessage(r, args)
		logger := getLogger(args)
		if nil != logger {
			logger.Error(message)
		} else {
			log.Println(message)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getMessage(r interface{}, args ...interface{}) (response string) {
	method := ""
	for _, item := range args {
		if s, ok := item.(string); ok {
			method = s
			break
		}
	}
	if len(method) > 0 {
		response = fmt.Sprintf("Error in application '%s' on '%s': %s", AppName, method, r)
	} else {
		response = fmt.Sprintf("[panic] Generic error in application '%s': %s", AppName, r)
	}
	return
}

func getLogger(args ...interface{}) (response qb_log.ILogger) {
	for _, item := range args {
		if logger, ok := item.(qb_log.ILogger); ok {
			response = logger
			break
		}
	}
	return
}
