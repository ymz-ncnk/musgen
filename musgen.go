//go:generate go run gen/dvar.go
package musgen

// MusGen is a code generator for the MUS format.
type MusGen interface {
	// Generate generates code from type description for the specified language.
	Generate(td TypeDesc, lang Lang) ([]byte, error)
}
