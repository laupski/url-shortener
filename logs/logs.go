package logs

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Format   string `toml:"format"`
	LogLevel string `toml:"log_level"`
	LogDest  string `toml:"log_dst"`
}

func SetLog(c Config) {

	switch c.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}

	switch c.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	if c.LogDest != "" {
		f, err := os.Create(c.LogDest)
		if err != nil {
			logrus.Fatal("create log file error: ", err.Error())
		}
		logrus.SetOutput(f)
	}
}
