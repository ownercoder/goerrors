package goerrors

import (
	"io"

	"github.com/sirupsen/logrus"
)

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func NewWriterHook(writer io.Writer, levels ...logrus.Level) logrus.Hook {
	return &WriterHook{
		Writer:    writer,
		LogLevels: levels,
	}
}

func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}
