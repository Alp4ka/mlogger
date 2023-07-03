package field

type integers interface {
	int | int64 | int32 | int16 | int8
}

type floats interface {
	float64 | float32
}
