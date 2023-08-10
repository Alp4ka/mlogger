package jsonsecurity

type Config struct {
	MaxDepth int
	Triggers map[string]TriggerOpts
}

type TriggerOpts struct {
	CaseSensitive bool
	MaskMethod    MaskerLabel
	ShouldAppear  bool

	original string
}
