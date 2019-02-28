package types

type TranscoderType int

const (
	BareMetal TranscoderType = iota
	Docker
)

type TranscoderRegisterInfo struct {
	Name string
	Port int
	Type TranscoderType
}
