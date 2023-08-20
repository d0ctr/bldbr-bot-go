package common

import (
	"fmt"
	standartLogger "log"
)

type LOGLEVEL byte

const (
	LEVELNOISE LOGLEVEL = 0
	LEVELDEBUG LOGLEVEL = 1
	LEVELINFO  LOGLEVEL = 2
	LEVELWARN  LOGLEVEL = 3
	LEVELERROR LOGLEVEL = 4
	LEVELFATAL LOGLEVEL = 5
)

func getLevelS(level LOGLEVEL) string {
	switch level {
	case LEVELINFO:
		return "INFO"
	case LEVELDEBUG:
		return "DEBUG"
	case LEVELNOISE:
		return "NOISE"
	case LEVELWARN:
		return "WARN"
	case LEVELERROR:
		return "ERROR"
	case LEVELFATAL:
		return "FATAL"
	default:
		return ""
	}
}

const prefixLine string = "%s %s: "

type loggerOptions struct {
	level  LOGLEVEL
	module string
}

type Logger struct {
	options loggerOptions
	logger  *standartLogger.Logger
}

func (logger *Logger) setPrefix(level LOGLEVEL) *standartLogger.Logger {
	logger.logger.SetPrefix(fmt.Sprintf(prefixLine, getLevelS(level), logger.options.module))
	return logger.logger
}

func (logger *Logger) LogF(level LOGLEVEL, line string, args ...any) {
	if level > logger.options.level {
		return
	}
	logger.setPrefix(level).Printf(line, args...)
}

func (logger *Logger) LogFatalF(line string, args ...any) {
	logger.setPrefix(LEVELFATAL).Fatalf(line, args...)
}

func (logger *Logger) Log(level LOGLEVEL, args ...any) {
	if level > logger.options.level {
		return
	}
	logger.setPrefix(level).Print(args...)
}

func (logger *Logger) LogFatal(args ...any) {
	logger.setPrefix(LEVELFATAL).Fatal(args...)
}

func CreateLogger(module string, level LOGLEVEL) *Logger {
	logger := Logger{
		logger: standartLogger.New(
			standartLogger.Writer(),
			"",
			standartLogger.Ldate|standartLogger.Ltime|standartLogger.Lmsgprefix,
		),
		options: loggerOptions{module: module, level: level},
	}

	return &logger
}

func (log *Logger) Child(module string, level LOGLEVEL) *Logger {
	module = log.options.module + "@" + module
	childLogger := CreateLogger(module, level)
	return childLogger
}
