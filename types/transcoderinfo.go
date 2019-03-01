package types

type TranscoderType int

const (
	BareMetal TranscoderType = iota
	Docker
)

type TranscoderInfo struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
	Type TranscoderType `json:"type"`
}