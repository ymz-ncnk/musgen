package musgen

const (
	// GoLang constant.
	GoLang Lang = iota
)

// Lang represents supported languages.
type Lang int

func (l Lang) String() string {
	return [...]string{"go"}[l]
}
