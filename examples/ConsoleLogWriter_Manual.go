package main

import (
	"time"

	l4g "log4go"
)

func main() {
	log := l4g.NewLogger()
	defer log.Close()
	log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}
