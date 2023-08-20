package logger

import (
	"fmt"
	standartLogger "log"
	"os"
	"strconv"
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

type LoggerOptions struct {
	Level  LOGLEVEL // max loglevel for the logger
	Module string   // max printable loglevel, may also be aquired from env var
}

type Logger struct {
	options *LoggerOptions
	logger  *standartLogger.Logger
}

func (logger *Logger) setPrefix(level LOGLEVEL) *standartLogger.Logger {
	logger.logger.SetPrefix(fmt.Sprintf(prefixLine, getLevelS(level), logger.options.module))
	return logger.logger
}

func (logger *Logger) LogF(level LOGLEVEL, line string, args ...any) {
	if level < logger.options.Level {
		return
	}
	logger.setPrefix(level).Printf(line, args...)
}

func (logger *Logger) LogFatalF(line string, args ...any) {
	logger.setPrefix(LEVELFATAL).Fatalf(line, args...)
}

func (logger *Logger) Log(level LOGLEVEL, args ...any) {
	if level < logger.options.Level {
		return
	}
	logger.setPrefix(level).Print(args...)
}

func (logger *Logger) LogFatal(args ...any) {
	logger.setPrefix(LEVELFATAL).Fatal(args...)
}

func CreateLogger(loggerOptions *LoggerOptions) *Logger {
	env_level := os.Getenv("LOGLEVEL")
	if loggerOptions.Level == 0 && len(env_level) != 0 {
		val, err := strconv.ParseInt(env_level, 10, 0)
		if err != nil {
			// doing nothing
		}
		loggerOptions.Level = LOGLEVEL(val)
	}
	logger := Logger{
		logger: standartLogger.New(
			standartLogger.Writer(),
			"",
			standartLogger.Ldate|standartLogger.Ltime|standartLogger.Lmsgprefix,
		),
		options: loggerOptions,
	}

	return &logger
}

func (Logger *Logger) Child(loggerOptions *LoggerOptions) *Logger {
	loggerOptions.Module = Logger.options.Module + "@" + loggerOptions.Module
	childLogger := CreateLogger(loggerOptions)
	return childLogger
}
