package simple

type SimpleModel struct {
	Some     int
	EvenMore int32 `rpc:"More"`
	Basic    string
	Value    string
}
