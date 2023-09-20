package field

type Option interface {
	Unpack() Field
}

type option struct {
	unpack func() Field
}

func (i option) Unpack() Field {
	return i.unpack()
}

func OptionCallerFunc(level ...int) Option {
	return option{
		unpack: func() Field { return CallerFunc(KeyCaller, level...) },
	}
}

func OptionSource(source string) Option {
	return option{
		unpack: func() Field { return String(KeySource, source) },
	}
}
