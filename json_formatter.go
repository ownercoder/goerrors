package goerrors

import "github.com/sirupsen/logrus"

type JSONFormatter struct {
	logrus.JSONFormatter
	IgnoreKeys []string
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	cleanEntry := filterEntryData(entry, f.IgnoreKeys)

	stack := getEntryStack(entry)
	if len(stack) > 0 {
		cleanEntry.Data[FormattedStackKey] = string(stack)
	}

	return f.JSONFormatter.Format(cleanEntry)
}
