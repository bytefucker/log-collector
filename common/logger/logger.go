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

	Instance.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	//Instance.SetReportCaller(true)
}

//设置日志级别
func EnableDebugLevel(b bool) {
	if b {
		Instance.SetLevel(logrus.DebugLevel)
	} else {
		Instance.SetLevel(logrus.InfoLevel)
	}
}
