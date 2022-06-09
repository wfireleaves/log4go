package log4go

type FieldType uint8

const (
	UnknownType FieldType = iota
	BoolType
	IntType
	Int32Type
	Uint32Type
	Int64Type
	Uint64Type
	Int8Type
	Uint8Type
	Float64Type
	Float32Type
	StringType
	InterfaceType
)

type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

func (f Field) AddTo(enc *jsonEncoder) {
	switch f.Type {
	case BoolType:
		enc.AddBool(f.Key, f.Interface.(bool))
	case IntType:
		enc.AddInt(f.Key, int(f.Integer))
	case Int32Type:
		enc.AddInt32(f.Key, int32(f.Integer))
	case Uint32Type:
		enc.AddUint32(f.Key, uint32(f.Integer))
	case Int64Type:
		enc.AddInt64(f.Key, f.Integer)
	case Uint64Type:
		enc.AddUint64(f.Key, uint64(f.Integer))
	case Int8Type:
		enc.AddInt8(f.Key, int8(f.Integer))
	case Uint8Type:
		enc.AddUint8(f.Key, int8(f.Integer))
	case Float32Type:
		enc.AddFloat32(f.Key, f.Interface.(float32))
	case Float64Type:
		enc.AddFloat64(f.Key, f.Interface.(float64))
	case StringType:
		enc.AddString(f.Key, f.String)
	case InterfaceType:
		enc.AddInterface(f.Key, f.Interface)
	}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Type: BoolType, Interface: value}
}

func Uint32(key string, value uint32) Field {
	return Field{Key: key, Type: Uint32Type, Integer: int64(value)}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Type: Int32Type, Integer: int64(value)}
}

func Uint8(key string, value uint8) Field {
	return Field{Key: key, Type: Uint8Type, Integer: int64(value)}
}

func Int8(key string, value int8) Field {
	return Field{Key: key, Type: Int8Type, Integer: int64(value)}
}

// Uint64 warning: max is int64
func Uint64(key string, value uint64) Field {
	return Field{Key: key, Type: Uint64Type, Integer: int64(value)}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Type: IntType, Integer: int64(value)}
}

func Float32(key string, value float32) Field {
	return Field{Key: key, Type: Float32Type, Interface: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Type: Float64Type, Interface: value}
}

func String(key string, value string) Field {
	return Field{Key: key, Type: StringType, String: value}
}

func Err(err error) Field {
	value := "nil"
	if err != nil {
		value = err.Error()
	}
	return Field{Key: "error", Type: StringType, String: value}
}

func Any(key string, value interface{}) Field {
	f := Field{
		Key: key,
	}
	switch value.(type) {
	case int:
		f.Integer = int64(value.(int))
		f.Type = IntType
	case uint8:
		f.Integer = int64(value.(uint8))
		f.Type = IntType
	case int8:
		f.Integer = int64(value.(int8))
		f.Type = IntType
	case uint32:
		f.Integer = int64(value.(uint32))
		f.Type = IntType
	case int32:
		f.Integer = int64(value.(int32))
		f.Type = IntType
	case uint64:
		f.Integer = int64(value.(uint64))
		f.Type = IntType
	case int64:
		f.Integer = int64(value.(int64))
		f.Type = IntType
	case string:
		f.String = value.(string)
		f.Type = StringType
	default:
		f.Type = InterfaceType
		f.Interface = value
	}
	return f
}
