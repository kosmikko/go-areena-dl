package download

import (
	"os"

	l "github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
)

var (
	logger l.Logger
)

func getLogger() l.Logger {
	if logger != nil {
		return logger
	}
	logger = golog.New(os.Stderr, getLevel(os.Getenv("LOG")))
	return logger
}

func getLevel(levelName string) (level l.Level) {
	switch levelName {
	case "debug":
		level = l.Debug
	case "info":
		level = l.Info
	case "notice":
		level = l.Notice
	case "warning":
		level = l.Warning
	case "error":
		level = l.Error
	case "critical":
		level = l.Critical
	case "alert":
		level = l.Alert
	case "emergency":
		level = l.Emergency
	case "none":
		level = l.None
	default:
		level = l.Info
	}

	return
}
