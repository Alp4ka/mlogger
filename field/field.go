package field

import (
	"github.com/Alp4ka/mlogger/jsonsecurity"
	"github.com/Alp4ka/mlogger/misc"
	"log/slog"
	"time"
)

var ContextKeyLogFields = struct{}{}

// Field stores key-value pairs in order to map them into log structure.
type Field struct {
	Attr slog.Attr
	Type Type

	err error
}

// Error returns an error field. It's used to store the error occurred while parsing, creating the field.
func (f Field) Error() error {
	return f.err
}

// Masked returns masked copy of Field using provided masker.
func (f Field) Masked(masker *jsonsecurity.Masker) Field {
	ff := Field{Attr: f.Attr, Type: f.Type, err: f.err}
	ff.Mask(masker)
	return ff
}

func (f Field) Key() string {
	return f.Attr.Key
}

func (f Field) Value() any {
	return f.Attr.Value.Any()
}

// Mask masks current Field using provided masker. In case the field has the type different to TypeJSONEscapeSecure -
// it will not make any changes.
func (f *Field) Mask(masker *jsonsecurity.Masker) {
	if f.Type != TypeJSONEscapeSecure {
		return
	}

	masked, err := masker.Mask(f.Attr.Value.String())
	if err == nil {
		f.Attr.Value = slog.StringValue(masked)
	} else {
		f.err = err
	}
}

// CUSTOM FIELDS

// Error used for storing errors. It displays as `{"error": "database error: timeout"}`
func Error(err error) Field {
	return ErrorNamed(KeyError, err)
}

// ErrorNamed used for storing errors. It's pretty similar to Error function, but uses key argument as it's key.
func ErrorNamed(key string, err error) Field {
	return Field{Attr: slog.String(key, err.Error()), Type: TypeError, err: nil}
}

// Int used for storing integer values.
func Int[T integers](key string, value T) Field {
	return Field{Attr: slog.Int64(key, int64(value)), Type: TypeInt, err: nil}
}

// Uint used for storing unsigned integer values.
func Uint[T unsignedIntegers](key string, value T) Field {
	return Field{Attr: slog.Uint64(key, uint64(value)), Type: TypeUint, err: nil}
}

// Float used for storing floating-point values.
func Float[T floats](key string, value T) Field {
	return Field{Attr: slog.Float64(key, float64(value)), Type: TypeFloat, err: nil}
}

// String used for storing string values.
func String(key string, value string) Field {
	return Field{Attr: slog.String(key, value), Type: TypeString, err: nil}
}

// Bool used for storing boolean values.
func Bool(key string, value bool) Field {
	return Field{Attr: slog.Bool(key, value), Type: TypeBool, err: nil}
}

// JSONEscape used for storing json string with escaping characters.
func JSONEscape(key string, value []byte) Field {
	return Field{Attr: slog.String(key, string(value)), Type: TypeJSONEscape, err: nil}
}

// JSONEscapeSecure the same as JSONEscape but also masks the key-value pairs specified in config.
// Example:
//
// Data = `{\"password\": \"qwerty123\", \"email\": \"example@example.com\"}`
//
// Using PASSWORD label for "password" and EMAIL for "email" we will reach the next result:
//
// Output = `{\"password\": \"*********\", \"email\": \"e******@example.com\"}`
func JSONEscapeSecure(key string, value []byte) Field {
	return Field{Attr: slog.String(key, string(value)), Type: TypeJSONEscapeSecure, err: nil}
}

// Any used for storing values of type any.
func Any(key string, value any) Field {
	return Field{Attr: slog.Any(key, value), Type: TypeAny, err: nil}
}

// CallerFunc creates field with name of caller function (using shift specified in level argument).
func CallerFunc(key string, level ...int) Field {
	const (
		defaultLevelShift = 1
	)

	var lvl int

	if len(level) == 0 {
		lvl = defaultLevelShift
	} else {
		lvl = level[0]
	}

	callerFuncName := misc.GetCallerWithLevel(lvl)
	return Field{Attr: slog.String(key, callerFuncName), Type: TypeCallerFunc, err: nil}
}

// Time stores timestamp.
func Time(key string, time time.Time) Field {
	return Field{Attr: slog.Time(key, time), Type: TypeTimestamp, err: nil}
}
