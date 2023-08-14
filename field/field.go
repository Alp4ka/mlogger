package field

import (
	"context"
	"github.com/Alp4ka/mlogger/jsonsecurity"
)

var ContextLogFields = struct{}{}

// Field stores key-value pairs in order to map them into log structure.
// TODO(Gorkovets Roman): Make Value field for al supported types in order to avoid inappropriate use of any.
type Field struct {
	Key   string
	Value any

	err error
}

// Error returns an error field. It's used to store the error occurred while parsing, creating the field.
func (f *Field) Error() error {
	return f.err
}

// Fields alias for []Field as it's more comfortable to use it this way
type Fields []Field

// Map maps field as key-value pairs in order to satisfy zerolog With() instructions.
func (fields Fields) Map() map[string]any {
	const failPostfix = "_FAIL"

	result := make(map[string]any)

	if fields == nil {
		return result
	}

	for _, field := range fields {
		result[field.Key] = field.Value

		if field.Error() != nil {
			result[field.Key+failPostfix] = field.Error()
		}
	}

	return result
}

// CUSTOM FIELDS

// Error used for storing errors. It displays as `{"error": "database error: timeout"}`
func Error(err error) Field {
	const key = "error"
	return Field{key, err, nil}
}

// Int used for storing integer values.
func Int[T integers](key string, value T) Field {
	return Field{key, int64(value), nil}
}

// Float used for storing floating-point values.
func Float[T floats](key string, value T) Field {
	return Field{key, float64(value), nil}
}

// String used for storing string values.
func String(key string, value string) Field {
	return Field{key, value, nil}
}

// Bool used for storing boolean values.
func Bool(key string, value bool) Field {
	return Field{key, value, nil}
}

// JSONEscape used for storing json string with escaping characters.
func JSONEscape(key string, value []byte) Field {
	return Field{key, value, nil}
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
	data, err := jsonsecurity.GlobalMasker().Mask(value)
	if err != nil {
		return Field{key, value, err}
	}
	return JSONEscape(key, data)
}

// Any used for storing values of type any.
func Any(key string, value any) Field {
	return Field{key, value, nil}
}

// FieldsFromCtx extract fields from context. If ctx is nil, use context.Background()
func FieldsFromCtx(ctx context.Context) Fields {
	if ctx == nil {
		ctx = context.Background()
	}
	fields, ok := ctx.Value(ContextLogFields).(Fields)
	if !ok {
		fields = make(Fields, 0, 0)
	}

	return fields
}

// WithContextFields appends fields to provided context.
func WithContextFields(ctx context.Context, fields ...Field) context.Context {
	return context.WithValue(ctx, ContextLogFields, append(FieldsFromCtx(ctx), fields...))
}
