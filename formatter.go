package goerrors

import "github.com/sirupsen/logrus"

type Formatter interface {
	Format(entry *logrus.Entry) ([]byte, error)
}

func filterEntryData(entry *logrus.Entry, ignoredKeys []string) *logrus.Entry {
	ignoredKeyMap := make(map[string]interface{})
	for _, k := range ignoredKeys {
		ignoredKeyMap[k] = nil
	}

	data := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		if _, ok := ignoredKeyMap[k]; ok {
			continue
		}

		data[k] = v
	}

	cleanEntry := *entry
	cleanEntry.Data = data

	return &cleanEntry
}
