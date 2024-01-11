package misc

import "log/slog"

// SlogReplaceAttr returns slog ReplaceAttr function for handler that substitutes mlogger Level in slog format.
func SlogReplaceAttr() func([]string, slog.Attr) slog.Attr {
	return func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := Level(a.Value.Any().(slog.Level))
			a.Value = slog.StringValue(level.String())
		}

		return a
	}
}

// SlogLevel simply casts mlogger Level to slog.Level.
func SlogLevel(level Level) slog.Level {
	return slog.Level(level)
}
