package field

import (
	"context"
	"github.com/Alp4ka/mlogger/jsonsecurity"
)

var ContextLogFields = struct{}{}

// Field TODO Сделать как в запе, чтобы не через рефлексию все работало в конечном итоге.
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

func Error(err error) Field {
	const key = "error"
	return Field{key, err, nil}
}

func Int[T integers](key string, value T) Field {
	return Field{key, int64(value), nil}
}

func Float[T floats](key string, value T) Field {
	return Field{key, float64(value), nil}
}

func String(key string, value string) Field {
	return Field{key, value, nil}
}

func Bool(key string, value bool) Field {
	return Field{key, value, nil}
}

func JSONEscape(key string, value []byte) Field {
	return Field{key, value, nil}
}

func JSONEscapeSecure(key string, value []byte) Field {
	data, err := jsonsecurity.GlobalMasker().Mask(value)
	if err != nil {
		return Field{key, nil, err}
	}
	return JSONEscape(key, data)
}

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

func WithContextFields(ctx context.Context, fields ...Field) context.Context {
	return context.WithValue(ctx, ContextLogFields, append(FieldsFromCtx(ctx), fields...))
}
