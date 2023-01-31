package logging

import (
	// External imports
	"github.com/sirupsen/logrus"
)

var (
	Log = logrus.New()
)

func init() {
	// init() is a reserved function in golang that executes
	// when the package is first imported.

	Log.Formatter = &logrus.TextFormatter{
		DisableLevelTruncation: true,
		FullTimestamp:          true,
	}
}
