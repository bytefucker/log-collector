package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Instance = logrus.New()

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Instance.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	Instance.SetLevel(logrus.DebugLevel)

	Instance.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	Instance.SetReportCaller(true)
}
