package qb_log

import "github.com/rskvp/qb-core/qb_utils"

type LogHelper struct {
}

var Log *LogHelper
var global *Logger

func init() {
	Log = new(LogHelper)
	global = NewLogger().
		OutFile(true).
		RotateEnable(true).
		SetLevel(InfoLevel)
}

func (instance *LogHelper) New(level, filename string) ILogger {
	logger := NewLogger()
	logger.OutFile(true)
	logger.SetLevel(level)
	logger.SetFilename(filename)

	return logger
}

func (instance *LogHelper) NewNoRotate(level, filename string) ILogger {
	logger := NewLogger()
	logger.OutFile(true)
	logger.SetLevel(level)
	logger.SetFilename(filename)
	logger.RotateEnable(false)
	
	// delete log file if any
	_ = qb_utils.IO.Remove(filename)

	return logger
}

func (instance *LogHelper) NewEmpty() *Logger {
	logger := NewLogger()
	return logger
}

func (instance *LogHelper) Logger() *Logger {
	return global
}
