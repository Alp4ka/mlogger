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
}

// Fields alias for []Field as it's more comfortable to use it this way
type Fields []Field

func (fields Fields) Map() map[string]any {
	result := make(map[string]any)

	if fields == nil {
		return result
	}

	for _, field := range fields {
		result[field.Key] = field.Value
	}

	return result
}

// CUSTOM FIELDS

func Error(err error) Field {
	const key = "error"
	return Field{key, err}
}

func Int[T integers](key string, value T) Field {
	return Field{key, int64(value)}
}

func Float[T floats](key string, value T) Field {
	return Field{key, float64(value)}
}

func String(key string, value string) Field {
	return Field{key, value}
}

func Bool(key string, value bool) Field {
	return Field{key, value}
}

func JSONEscape(key string, value []byte) Field {
	return Field{key, value}
}

func JSONEscapeSecure(key string, value []byte) Field {
	// TODO(Gorkovets Roman): Crutch
	const failPostfix = "_FAIL"

	data, err := jsonsecurity.GlobalMasker().Mask(value)
	if err != nil {
		return Field{key + failPostfix, err.Error()}
	}
	return JSONEscape(key, data)
}

func Any(key string, value any) Field {
	return Field{key, value}
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
