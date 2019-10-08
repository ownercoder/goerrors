package goerrors

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Options struct {
	SentryDSN        string
	SentrySyncMode   bool
	SentrySkipFrames int
}

type Level string

const (
	PanicLevel   Level = "panic"
	FatalLevel   Level = "fatal"
	ErrorLevel   Level = "error"
	WarningLevel Level = "warning"
	InfoLevel    Level = "info"
	DebugLevel   Level = "debug"
	TraceLevel   Level = "trace"
)

var instance atomic.Value
var logrusLevelMap = map[Level]logrus.Level{
	PanicLevel:   logrus.PanicLevel,
	FatalLevel:   logrus.FatalLevel,
	ErrorLevel:   logrus.ErrorLevel,
	WarningLevel: logrus.WarnLevel,
	InfoLevel:    logrus.InfoLevel,
	DebugLevel:   logrus.DebugLevel,
	TraceLevel:   logrus.TraceLevel,
}

const (
	FormattedStackKey = "stack"
	HTTPRequestKey    = "http_request"
)

type Extra struct {
	Request     *http.Request
	Fingerprint []string
}

func newLogger(level Level, formatter Formatter, options *Options) (*logrus.Logger, error) {
	logger := logrus.New()

	if options == nil {
		options = &Options{}
	}

	logrusLevel, ok := logrusLevelMap[level]
	if !ok {
		return nil, errors.Errorf("invalid error level: %s", level)
	}

	logger.SetLevel(logrusLevel)
	logger.SetFormatter(formatter)
	logger.SetOutput(ioutil.Discard)

	errLevels := []logrus.Level{
		logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel,
	}
	infoLevels := []logrus.Level{
		logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel,
	}

	logger.AddHook(NewWriterHook(os.Stderr, errLevels...))
	logger.AddHook(NewWriterHook(os.Stdout, infoLevels...))

	sentryOptions := sentry.ClientOptions{
		Dsn: options.SentryDSN,
	}
	if options.SentrySyncMode {
		syncTransport := sentry.NewHTTPSyncTransport()
		syncTransport.Timeout = 3 * time.Second
		sentryOptions.Transport = syncTransport
	}

	sentryHook, err := NewSentryHook(sentryOptions, options.SentrySkipFrames, errLevels...)
	if err != nil {
		return nil, err
	}
	logger.AddHook(sentryHook)

	return logger, nil
}

func InitLog(level Level, formatter Formatter, options *Options) error {
	logger, err := newLogger(level, formatter, options)
	if err != nil {
		return err
	}

	instance.Store(logger)

	return nil
}

func Log() *logrus.Logger {
	if logger, ok := instance.Load().(*logrus.Logger); ok && logger != nil {
		return logger
	}

	return logrus.StandardLogger()
}
