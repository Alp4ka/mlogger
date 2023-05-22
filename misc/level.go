package misc

type Level int

const (
	LevelNone Level = iota
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
	LevelError: "Error",
	LevelFatal: "Fatal",
	LevelPanic: "Panic",
	LevelWarn:  "Warn",
}

func (l Level) String() string {
	if val, ok := _lvlToString[l]; ok {
		return val
	}
	return _lvlToString[LevelNone]
}
