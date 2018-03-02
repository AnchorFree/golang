// Package jsonlog implements a simple JSON logging package.
// It defines a Logger interface, and provides StdLogger
// type that satisfies that interface. It also provides
// a MockLogWriter type which might come handy for tests.
//
// A typical log message will have the following fields:
// { "timestamp":1519317784, "level":"debug", "message":"log message",
//   "err":"error message", "app":"appName" }
//
// App and timestamp are optional fields that can be
// omitted, see Init method. Level indicates the log level
// of the message. You can provide additional context
// to a log message, Debug, Info, Error and Fatal methods
// all accept optional context argument. For example:
//
// -- code --
//	type ctx struct {
//		ActionID string `json:"action_id"`
//		Duration int    `json:"duration,omitempty"`
//	}
//	log := jsonlog.Init("myApp", false, true, nil)
//	log.Info("log event", &ctx{ActionID: "id1", Duration: 10})
// -- code --
//
// Will produce the following log message:
// {"level":"info","app":"myApp","context":[{"action_id":"id1","duration":10}],"message":"log event"}
package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

// Logger is the interface that wraps basic logging methods.
type Logger interface {
	Init(appName string, debug, noTimestamp bool, w io.Writer)
	Info(msg string, context ...interface{})
	Debug(msg string, context ...interface{})
	Error(msg string, err error, context ...interface{})
	Fatal(msg string, err error, context ...interface{})
	Stats() (debug, info, errors, fatal int)
	Flush()
}

// StdLogger is an implementation of the Logger interface.
// Field "stats" is a map of counters for different logging levels,
// e.g. when you use Info() method, it will increase
// stats["info"] counter.
type StdLogger struct {
	stderr zerolog.Logger
	stats  map[string]int
}

// Init sets the defaults and overrides them with user settings.
// If appName is provided, then additional field { "app": "appName" }
// will be attached to all emitted log messages. You can provide
// an empty string as an appName if you don't want that.
// If you provide nil instead of io.Writer, os.Stderr will be used.
// A typical invocation, if you want defaults, would look like:
// Init("", false, false, nil)
func (L *StdLogger) Init(appName string, debug, noTimestamp bool, w io.Writer) {

	zerolog.TimeFieldFormat = ""
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if w == nil {
		w = os.Stderr
	}
	if noTimestamp {
		L.stderr = zerolog.New(w)
	} else {
		L.stderr = zerolog.New(w).With().Timestamp().Logger()
	}

	if appName != "" {
		L.stderr = L.stderr.With().Str("app", appName).Logger()
	}
	L.stats = map[string]int{"debug": 0, "info": 0, "error": 0, "fatal": 0}

}

func withContext(L zerolog.Logger, context ...interface{}) zerolog.Logger {

	if context != nil {
		return L.With().Interface("context", context).Logger()
	}
	return L

}

// Info writes a log msg with info log level.
func (L *StdLogger) Info(msg string, context ...interface{}) {

	log := withContext(L.stderr, context...)
	log.Info().Msg(msg)
	L.stats["info"]++

}

// Error writes a log msg with error log level.
func (L *StdLogger) Error(msg string, err error, context ...interface{}) {

	log := withContext(L.stderr, context...)
	log.Error().Err(err).Msg(msg)
	L.stats["error"]++

}

// Fatal writes a log msg with fatal log level and does os.Exit(1) after that.
func (L *StdLogger) Fatal(msg string, err error, context ...interface{}) {

	log := withContext(L.stderr, context...)
	log.Fatal().Err(err).Msg(msg)
	L.stats["fatal"]++

}

// Debug writes a log msg with debug log level. If debug log level
// was not set with Init method globally, then all log messages with
// debug level will get silently discarded.
func (L *StdLogger) Debug(msg string, context ...interface{}) {

	log := withContext(L.stderr, context...)
	log.Debug().Msg(msg)
	L.stats["debug"]++

}

// Flush zeros stats counters
func (L *StdLogger) Flush() {
	L.stats = map[string]int{"debug": 0, "info": 0, "error": 0, "fatal": 0}
}

// Stats returns the number of log messages that have been
// emitted so far for debug, info, error and fatal levels.
func (L *StdLogger) Stats() (debug, info, errors, fatal int) {

	debug = L.stats["debug"]
	info = L.stats["info"]
	errors = L.stats["error"]
	fatal = L.stats["fatal"]
	return

}

// MockLogWriter is an implementation of io.Writer interface,
// which can be used in testing.
type MockLogWriter struct {
	Logs []string
}

func (m *MockLogWriter) Write(p []byte) (n int, err error) {

	m.Logs = append(m.Logs, string(p))
	return len(p), nil

}
