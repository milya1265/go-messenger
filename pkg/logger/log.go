package logger

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
)

type Logger struct {
	*logrus.Entry
}

var Log Logger

func (s *Logger) ExtraFields(fields map[string]interface{}) *Logger {
	return &Logger{
		s.WithFields(fields)}
}

var instance Logger
var once sync.Once

func GetLogger(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}

		l := logrus.New()
		l.SetReportCaller(true)

		l.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
			},
			ForceColors:   true,
			FullTimestamp: true,
			PadLevelText:  true,
		})
		l.SetOutput(colorable.NewColorableStdout())

		l.SetOutput(os.Stdout)
		l.SetLevel(logrusLevel)

		instance = Logger{logrus.NewEntry(l)}
	})

	return instance
}
