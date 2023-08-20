package main

import (
	standardLogger "log"
)

func log(args ...any) {
	standardLogger.Default().Print(args...)
}

func logFatal(args ...any) {
	standardLogger.Default().Fatal(args...)
}

type loggerArgs struct {
	level  string
	module string
}

type formattedLogger struct {
	_log *standardLogger.Logger
}

func (log formattedLogger) log(line string, v ...any) {
	log._log.Printf(line, v...)
}

func (log formattedLogger) logFatal(line string, v ...any) {
	log._log.Fatalf(line, v...)
}

type logger struct {
	f        formattedLogger
	_options loggerArgs
	_log     *standardLogger.Logger
}

func (log logger) log(args ...any) {
	log._log.Print(args...)
}

func (log logger) logFatal(args ...any) {
	log._log.Fatal(args...)
}

func Logger(options loggerArgs) *logger {
	logger := logger{
		_log: standardLogger.New(
			standardLogger.Writer(),
			options.prefix+" ",
			standardLogger.Ldate|standardLogger.Ltime|standardLogger.Lmsgprefix,
		),
		_options: options,
	}
	logger.f._log = logger._log
	return &logger
}

func (log *logger) child(options loggerArgs) *logger {
	childLogger := Logger(loggerArgs{
		prefix: log._options.prefix + "@" + options.prefix,
	})
	return childLogger
}
