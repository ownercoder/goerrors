package goerrors

import (
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var sentryLevelMap = map[logrus.Level]sentry.Level{
	logrus.PanicLevel: sentry.LevelFatal,
	logrus.FatalLevel: sentry.LevelFatal,
	logrus.ErrorLevel: sentry.LevelError,
	logrus.WarnLevel:  sentry.LevelWarning,
	logrus.InfoLevel:  sentry.LevelInfo,
	logrus.DebugLevel: sentry.LevelDebug,
	logrus.TraceLevel: sentry.LevelDebug,
}

type SentryHook struct {
	skipFrames int
	client     *sentry.Client
	levels     []logrus.Level
}

func NewSentryHook(options sentry.ClientOptions, skipFrames int, levels ...logrus.Level) (logrus.Hook, error) {
	client, err := sentry.NewClient(options)
	if err != nil {
		return nil, err
	}

	hook := SentryHook{
		skipFrames: skipFrames,
		client:     client,
		levels:     levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook, nil
}

func (hook *SentryHook) Fire(entry *logrus.Entry) error {
	var exceptions []sentry.Exception

	if err, ok := entry.Data[logrus.ErrorKey].(error); ok && err != nil {
		exceptions = append(exceptions, sentry.Exception{
			Type:       entry.Message,
			Value:      err.Error(),
			Stacktrace: hook.prepareStacktrace(err),
		})
	}

	event := sentry.Event{
		Level:       sentryLevelMap[entry.Level],
		Message:     entry.Message,
		Fingerprint: getEntryFingerprint(entry),
		Tags:        getEntryTags(entry),
		Extra:       map[string]interface{}(entry.Data),
		Exception:   exceptions,
	}

	if request := getEntryHTTPRequest(entry); request != nil {
		event.Request = sentry.Request{}.FromHTTPRequest(request)
	}

	hub := sentry.CurrentHub()
	hook.client.CaptureEvent(&event, nil, hub.Scope())

	return nil
}

func (hook *SentryHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *SentryHook) prepareStacktrace(err error) *sentry.Stacktrace {
	stacktrace := sentry.ExtractStacktrace(err)
	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}

	if hook.skipFrames < 1 {
		return stacktrace
	}

	framesCount := len(stacktrace.Frames)
	if framesCount <= hook.skipFrames {
		return nil
	}

	stacktrace.Frames = stacktrace.Frames[:framesCount-hook.skipFrames]

	return stacktrace
}
