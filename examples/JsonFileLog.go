package main

import l4g "log4go"

func main() {
	// Get a new logger instance
	log := l4g.NewLogger()

	// Create a default logger that is logging messages of FINE or higher
	log.AddFilter("file", l4g.FINE, l4g.NewFileLogWriter("test_json", false))
	log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	defer log.Close()

	log.Errorjson("test json", l4g.Int32("int32", 111), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
}
