package log

import (
	"fmt"
	"os"

	logrus "github.com/Sirupsen/logrus"
	cfg "github.com/mccraymt/ms-black-history/config"
	loggly "github.com/sebest/logrusly"
)

func init() {
	fmt.Println("Initializing logger")
	// Output to stderr instead of stdout, could also be a file.
	logrus.SetOutput(os.Stderr)
	var logLevel logrus.Level
	// configure logging settings based on environment
	fmt.Printf("Setting log level to %v \n", cfg.Config.LogLevel)
	switch cfg.Config.LogLevel {
	case "debug":
		logLevel = logrus.DebugLevel
		break
	case "info":
		logLevel = logrus.InfoLevel
		break
	case "warning":
		logLevel = logrus.WarnLevel
		break
	case "error":
		logLevel = logrus.ErrorLevel
		break
	case "fatal":
		logLevel = logrus.FatalLevel
		break
	case "panic":
		logLevel = logrus.PanicLevel
	}

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// configure logrus log level
	logrus.SetLevel(logLevel)

	// configure loggly
	hostName, _ := os.Hostname()
	logglyHook := loggly.NewLogglyHook(cfg.Config.LogglyKey, hostName, logLevel, cfg.Config.Environment, "ms-geo-data")
	logrus.Println("Adding Loggly logging hook")
	logrus.AddHook(logglyHook)
}
