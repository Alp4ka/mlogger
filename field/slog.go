package field

import "log/slog"

func UnpackFieldsToSlogAttrs(fields Fields) []slog.Attr {
	attrs := make([]slog.Attr, len(fields))

	for i := range fields {
		attrs[i] = fields[i].Attr
	}

	return attrs
}
