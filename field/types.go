package field

type integers interface {
	int | int64 | int32 | int16 | int8
}

type floats interface {
	float64 | float32
}

type unsignedIntegers interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Type uint8

const (
	TypeNone Type = iota
	TypeError
	TypeInt
	TypeUint
	TypeFloat
	TypeString
	TypeStringer
	TypeBool
	TypeJSONEscape
	TypeJSONEscapeSecure
	TypeAny
	TypeCallerFunc
	TypeTimestamp
)
