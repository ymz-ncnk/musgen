package musgen

// TypeDesc represents a type description. It is a pipeline for the
// struct.go.tmpl.
type TypeDesc struct {
	Package string
	Name    string
	Unsafe  bool
	Suffix  string
	Fields  []FieldDesc
}

// FieldDesc represents a field description. It is a part of the TypeDesc.
type FieldDesc struct {
	Name          string
	Type          string
	MaxLength     int
	Alias         string
	Validator     string
	Encoding      string
	KeyValidator  string
	KeyEncoding   string
	ElemValidator string
	ElemEncoding  string
}
