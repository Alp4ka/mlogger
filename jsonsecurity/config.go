package jsonsecurity

type Config struct {
	Triggers map[string]TriggerOpts
}

type TriggerOpts struct {
	CaseSensitive bool
	MaskMethod    MaskerLabel
	ShouldAppear  bool
}
