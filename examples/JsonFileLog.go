package main

import (
	"github.com/wfireleaves/log4go"
	l4g "github.com/wfireleaves/log4go"
)

func main() {
	// Get a new logger instance
	log4go.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("test_json", false))
	// Create a default logger that is logging messages of FINE or higher
	defer log4go.Close()

	log4go.Error("test json1", l4g.Int32("int32", 1), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
	log4go.Error("test json2", l4g.Int32("int32", 11), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
	log4go.Error("test json3", l4g.Int32("int32", 111), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
	log4go.Error("test json4", l4g.Int32("int32", 1111), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
	log4go.Error("test json5", l4g.Int32("int32", 11111), l4g.String("string", "aaaa"), l4g.Bool("bool", true))
	log4go.Error("test json float32", l4g.Float32("float32", 12.3456))
	log4go.Error("test json float64", l4g.Float64("float64", 12345123152.345124))
	log4go.Error("test json uint8", l4g.Uint8("uint8", 8))
	log4go.Error("test json bool", l4g.Bool("true", false))
}
