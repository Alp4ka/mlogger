package misc

import "log/slog"

func SlogReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := Level(a.Value.Any().(slog.Level))
		a.Value = slog.StringValue(level.String())
	}

	return a
}

func SlogLevel(level Level) slog.Level {
	return slog.Level(level)
}
