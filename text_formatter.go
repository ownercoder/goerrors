package goerrors

import "github.com/sirupsen/logrus"

type TextFormatter struct {
	logrus.TextFormatter
	IgnoreKeys []string
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	stack := getEntryStack(entry)

	text, err := f.TextFormatter.Format(filterEntryData(entry, f.IgnoreKeys))
	if err != nil {
		return nil, err
	}

	if len(stack) > 0 {
		text = append(text, "Trace:\n"...)
		text = append(text, stack...)
	}

	return text, nil
}
