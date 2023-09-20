package field

import (
	"context"
	"github.com/Alp4ka/mlogger/jsonsecurity"
)

// Fields alias for []Field as it's more comfortable to use it this way
type Fields []Field

// Prepare maps field as key-value pairs in order to satisfy zerolog With() instructions.
func (fields Fields) Prepare(masker *jsonsecurity.Masker) Fields {
	// Build key string of fields with not nil err field inside.
	// Example:
	//
	// my_failed_field -> my_failed_field_FAIL
	buildErrorField := func(f Field) string {
		return f.Attr.Key + "_FAIL"
	}

	if fields == nil {
		return make(Fields, 0)
	}

	// Since every field can have an error inside it's structure we suppose that no more than 2*N fields will be in
	// resulting slice.
	ret := make(Fields, 0, 2*len(fields))

	// Go through fields slice and append the same field with suffix _FAIL when not nil err stored in Field structure.
	// Also prepares other fields basing on their types.
	for _, field := range fields {
		ret = append(ret, field.Masked(masker))
		fieldErr := field.Error()
		if fieldErr != nil {
			ret = append(ret, ErrorNamed(buildErrorField(field), fieldErr))
		}
	}

	return ret
}

// WithOptions copies fields and appends this copy with specified options calling Option.Unpack() method.
func (fields Fields) WithOptions(options ...Option) Fields {
	fieldsCopy := make(Fields, 0, len(fields)+len(options))

	fieldsCopy = append(fieldsCopy, fields...)
	for _, o := range options {
		fieldsCopy = append(fieldsCopy, o.Unpack())
	}

	return fieldsCopy
}

// FieldsFromCtx extract fields from context. If ctx is nil, use context.Background()
// Deprecated: Use FieldsFromContext instead.
func FieldsFromCtx(ctx context.Context) Fields {
	return FieldsFromContext(ctx)
}

// FieldsFromContext extract fields from context. If ctx is nil, use context.Background()
func FieldsFromContext(ctx context.Context) Fields {
	if ctx == nil {
		ctx = context.Background()
	}

	fields, ok := ctx.Value(ContextKeyLogFields).(Fields)
	if !ok {
		fields = make(Fields, 0, 0)
	}

	return fields
}

// WithContextFields appends fields to provided context.
func WithContextFields(ctx context.Context, fields ...Field) context.Context {
	return context.WithValue(ctx, ContextKeyLogFields, append(FieldsFromContext(ctx), fields...))
}
