package musgen

// MusGen is a code generator for the MUS format.
type MusGen interface {
	// Generate generates code for the specified language.
	Generate(td TypeDesc, lang Lang) ([]byte, error)
}
