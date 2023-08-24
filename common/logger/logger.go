package logger

import (
	"fmt"
	standartLogger "log"
	"os"
	"strconv"
	"strings"
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
	case LEVELNOISE:
		return "NOISE"
	case LEVELDEBUG:
		return "DEBUG"
	case LEVELINFO:
		return "INFO"
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

func getLevelColoring(level LOGLEVEL) string {
	switch level {
	case LEVELNOISE:
		return "\033[1;35m%5s\033[0m"
	case LEVELDEBUG:
		return "\033[1;34m%5s\033[0m"
	case LEVELINFO:
		return "\033[1;32m%5s\033[0m"
	case LEVELWARN:
		return "\033[5m%5s\033[0m"
	case LEVELERROR:
		return "\033[1;31m%5s\033[0m"
	case LEVELFATAL:
		return "\033[31m%5s\033[0m"
	default:
		return "%5s"
	}
}

const prefixLine string = "[%5s] %s -> "

type LoggerOptions struct {
	Level     LOGLEVEL // max loglevel for the logger
	Module    string   // max printable loglevel, may also be aquired from env var
	IsColored bool
}

type Logger struct {
	options *LoggerOptions
	logger  *standartLogger.Logger
	F       FLogger // logger, which print methods receive formatted string as paramters
}

func (logger *Logger) setPrefix(level LOGLEVEL) *standartLogger.Logger {
	prefixLine := prefixLine
	if logger.options.IsColored {
		prefixLine = fmt.Sprintf(prefixLine, getLevelColoring(level), "%s")
	}

	logger.logger.SetPrefix(fmt.Sprintf(prefixLine, getLevelS(level), logger.options.Module))
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
		val, _ := strconv.ParseInt(env_level, 10, 0)
		loggerOptions.Level = LOGLEVEL(val)
	}
	loggerOptions.Module = strings.ReplaceAll(loggerOptions.Module, " ", ".")
	logger := &Logger{
		logger: standartLogger.New(
			standartLogger.Writer(),
			"",
			standartLogger.Ldate|standartLogger.Ltime|standartLogger.Lmsgprefix,
		),
		options: loggerOptions,
	}
	logger.F = FLogger{
		logger: logger,
	}

	return logger
}

type ChildLoggerOptions struct {
	Level     any // max loglevel for the logger
	Module    any // max printable loglevel, may also be aquired from env var
	IsColored any
}

func (logger *Logger) Child(newOptions *ChildLoggerOptions) *Logger {
	loggerOptions := *logger.options
	if newOptions.Module != nil {
		module, ok := newOptions.Module.(string)
		if ok {
			loggerOptions.Module = loggerOptions.Module + "@" + module
		}
	}
	if newOptions.Level != nil {
		level, ok := newOptions.Level.(LOGLEVEL)
		if ok {
			loggerOptions.Level = level
		}
	}
	if newOptions.IsColored != nil {
		isColored, ok := newOptions.IsColored.(bool)
		if ok {
			loggerOptions.IsColored = isColored
		}
	}

	childLogger := CreateLogger(&loggerOptions)
	return childLogger
}

func (logger *Logger) NOISE(args ...any) {
	logger.Log(LEVELNOISE, args...)
}

func (logger *Logger) DEBUG(args ...any) {
	logger.Log(LEVELDEBUG, args...)
}

func (logger *Logger) INFO(args ...any) {
	logger.Log(LEVELINFO, args...)
}

func (logger *Logger) WARN(args ...any) {
	logger.Log(LEVELWARN, args...)
}

func (logger *Logger) ERROR(args ...any) {
	logger.Log(LEVELERROR, args...)
}

func (logger *Logger) FATAL(args ...any) {
	logger.LogFatal(args...)
}

type FLogger struct {
	logger *Logger
}

func (f *FLogger) NOISE(line string, args ...any) {
	f.logger.LogF(LEVELNOISE, line, args...)
}

func (f *FLogger) DEBUG(line string, args ...any) {
	f.logger.LogF(LEVELDEBUG, line, args...)
}

func (f *FLogger) INFO(line string, args ...any) {
	f.logger.LogF(LEVELINFO, line, args...)
}

func (f *FLogger) WARN(line string, args ...any) {
	f.logger.LogF(LEVELWARN, line, args...)
}

func (f *FLogger) ERROR(line string, args ...any) {
	f.logger.LogF(LEVELERROR, line, args...)
}

func (f *FLogger) FATAL(line string, args ...any) {
	f.logger.LogFatalF(line, args...)
}

var RootLogger = CreateLogger(&LoggerOptions{
	Module:    "root",
	IsColored: false,
})

func SetRootLogger(logger *Logger) {
	RootLogger = logger
}
