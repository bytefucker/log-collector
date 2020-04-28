package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

var Instance *Logger

func init() {
	if Instance == nil {
		Instance = New()
	}
	Instance.SetOutput(os.Stdout)
	Instance.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Instance.SetMode(logrus.InfoLevel)
	//Instance.ReportCaller = true
}

func (l *Logger) SetMode(level logrus.Level) {
	l.SetLevel(level)
}

func New() *Logger {
	return &Logger{logrus.New()}
}
