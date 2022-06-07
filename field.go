package log4go

type FieldType uint8

const (
	UnknownType FieldType = iota
	BoolType
	Int32Type
	StringType
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
	case Int32Type:
		enc.AddInt32(f.Key, int32(f.Integer))
	case StringType:
		enc.AddString(f.Key, f.String)
	}
}
