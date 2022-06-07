package log4go

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

const Hex = "0123456789ABCDEF"

type jsonEncoder struct {
	buf  []byte
	left bool
}

func NewJsonEncoder() *jsonEncoder {
	return &jsonEncoder{
		buf: make([]byte, 0, 100),
	}
}

func (enc *jsonEncoder) EncodeJson(record *LogRecord) string {
	enc.appendByte('{')
	enc.appendString(`"time":`)
	enc.appendString(fmt.Sprintf("\"%04d-%02d-%02d %02d:%02d:%02d.%05d\"",
		record.Created.Year(), record.Created.Month(), record.Created.Day(),
		record.Created.Hour(), record.Created.Minute(), record.Created.Second(), record.Created.Nanosecond()/10000))
	enc.appendString(`,"message":`)
	enc.appendString(`"` + record.Message + `"`)
	enc.appendString(`,"level":`)
	enc.appendString(`"` + record.Level.String() + `"`)
	enc.appendString(`,"file":`)
	enc.appendString(`"` + record.Source + `"`)
	for _, f := range record.Fields {
		f.AddTo(enc)
	}
	enc.appendByte('}')
	enc.appendByte('\n')
	return string(enc.buf)
}

func (enc *jsonEncoder) EncodeString(record *LogRecord) string {
	if len(record.Fields) <= 0 {
		return record.Message
	}
	enc.appendString(record.Message)
	args := make([]interface{}, 0, len(record.Fields))
	for _, f := range record.Fields {
		enc.appendString(f.Key + ":")
		switch f.Type {
		case Int32Type:
			enc.appendString("%d")
			args = append(args, f.Integer)
		case StringType:
			enc.appendString("%s")
			args = append(args, f.String)
		default:
			enc.appendString("%v")
			args = append(args, f.Interface)
		}
		enc.appendByte(' ')
	}
	format := string(enc.buf)
	format = fmt.Sprintf(format, args...)
	return format
}

func (enc *jsonEncoder) AddInt32(key string, value int32) {
	enc.AppendLeft()
	enc.safeAddString(key)
	enc.appendString(`": `)
	enc.appendString(strconv.Itoa(int(value)))
}

func (enc *jsonEncoder) AddBool(key string, value bool) {
	enc.AppendLeft()
	enc.safeAddString(key)
	enc.appendString(`": `)
	enc.buf = strconv.AppendBool(enc.buf, value)
}

func (enc *jsonEncoder) AddString(key, value string) {
	enc.AppendLeft()
	enc.safeAddString(key)
	enc.appendString(`": "`)
	enc.safeAddString(value)
	enc.appendByte('"')
}

func (enc *jsonEncoder) AppendLeft() {
	if enc.left == true {
		enc.appendString(`"`)
		enc.left = false
	} else {
		enc.appendString(`,"`)
	}
}

func (enc *jsonEncoder) appendString(str string) {
	enc.buf = append(enc.buf, str...)
}

func (enc *jsonEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf = append(enc.buf, s[i:i+size]...)
		i += size
	}
}

func (enc *jsonEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.appendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.appendByte('\\')
		enc.appendByte(b)
	case '\n':
		enc.appendByte('\\')
		enc.appendByte('n')
	case '\r':
		enc.appendByte('\\')
		enc.appendByte('r')
	case '\t':
		enc.appendByte('\\')
		enc.appendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf = append(enc.buf, `\u00`...)
		enc.buf = append(enc.buf, Hex[b>>4])
		enc.buf = append(enc.buf, Hex[b&0xF])
	}
	return true
}

func (enc *jsonEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf = append(enc.buf, `\ufffd`...)
		return true
	}
	return false
}

func (enc *jsonEncoder) appendByte(b byte) {
	enc.buf = append(enc.buf, b)
}
