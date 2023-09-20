package misc

type Level int

const (
	LevelNone = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var _lvlToString = map[Level]string{
	LevelNone:  "None",
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelPanic: "Panic",
	LevelFatal: "Fatal",
}

func (l Level) String() string {
	if val, ok := _lvlToString[l]; ok {
		return val
	}
	return _lvlToString[LevelNone]
}

func (l Level) LessThan(level Level) bool {
	return l < level
}

func (l Level) LessOrEqualThan(level Level) bool {
	return l <= level
}

func (l Level) BiggerThan(level Level) bool {
	return l > level
}

func (l Level) BiggerOrEqualThan(level Level) bool {
	return l >= level
}

func (l Level) EqualTo(level Level) bool {
	return l == level
}
