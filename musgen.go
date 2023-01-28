//go:generate go run metagen/templates_var.go
package musgen

import (
	"bytes"
	"regexp"
	"strconv"
	tmplmod "text/template"
)

// Two main templates are templates/alias.go.tmpl and templates/struct.go.tmpl.
// After any changes in one of the template file from templates/ directory, you
// should run '$ go generate'.

// UintWithSystemSize uint type with system size.
var UintWithSystemSize = "uint" + strconv.Itoa(strconv.IntSize)

const (
	MaxVarintLength64 = 10
	MaxVarintLength32 = 5
	MaxVarintLength16 = 3
	MaxVarintLength8  = 1
	RawEncoding       = "raw"
)

// New returns a new MusGen.
func New() (MusGen, error) {
	baseTmpl := tmplmod.New("base")
	registerFuncs(baseTmpl)
	err := registerTemplates(baseTmpl)
	if err != nil {
		return MusGen{}, err
	}
	return MusGen{baseTmpl}, err
}

// MusGen is a code generator for the MUS format.
type MusGen struct {
	baseTmpl *tmplmod.Template
}

// Generate generates a code for the specified language.
func (gen MusGen) Generate(td TypeDesc, lang Lang) ([]byte, error) {
	err := gen.validateTypeDesc(td)
	if err != nil {
		return nil, err
	}
	switch lang {
	case GoLang:
		if Alias(td) {
			return gen.generate(td, makeAliasFileName(lang))
		}
		return gen.generate(td, makeStructFileName(lang))
	default:
		return nil, UnsupportedLangError{lang}
	}
}

func (gen MusGen) generate(td TypeDesc, tmplFile string) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := gen.baseTmpl.ExecuteTemplate(buf, tmplFile, td)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (gen MusGen) validateTypeDesc(td TypeDesc) error {
	return gen.validateEncoding(td)
}

func (gen MusGen) validateEncoding(td TypeDesc) (err error) {
	for _, field := range td.Fields {
		if field.Encoding != "" {
			err = gen.validateFieldEncoding(field.Type, field.Encoding)
			if err != nil {
				return err
			}
		}
		if field.ElemEncoding != "" {
			err = gen.validateElemEncoding(field.Type, field.ElemEncoding)
			if err != nil {
				return err
			}
		}
		if field.KeyEncoding != "" {
			err = gen.validateKeyEncoding(field.Type, field.KeyEncoding)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gen MusGen) validateFieldEncoding(fieldType, encoding string) error {
	if !gen.supportEncoding(fieldType, encoding) {
		return ErrUnsupportedEncoding
	}
	return nil
}

func (gen MusGen) validateElemEncoding(fieldType, encoding string) error {
	var elemType string
	if at := ParseArrayType(fieldType); at.Valid {
		elemType = at.Type
	} else if st := ParseSliceType(fieldType); st.Valid {
		elemType = st.Type
	} else if mt := ParseMapType(fieldType); mt.Valid {
		elemType = mt.Value
	} else {
		return ErrUnsupportedElemEncoding
	}
	if !gen.supportEncoding(elemType, encoding) {
		return ErrUnsupportedElemEncoding
	}
	return nil
}

func (gen MusGen) validateKeyEncoding(fieldType, encoding string) error {
	var keyType string
	if mt := ParseMapType(fieldType); mt.Valid {
		keyType = mt.Key
	} else {
		return ErrUnsupportedKeyEncoding
	}
	if !gen.supportEncoding(keyType, encoding) {
		return ErrUnsupportedKeyEncoding
	}
	return nil
}

func (gen MusGen) supportEncoding(baseTmpl string, encoding string) bool {
	if encoding == RawEncoding && gen.supportRawEncoding(baseTmpl) {
		return true
	}
	return false
}

func (gen MusGen) supportRawEncoding(baseTmpl string) bool {
	re := regexp.MustCompile(`^\**(?:(?:uint|int)(?:64|32|16|8|)|float(?:64|32))$`)
	return re.MatchString(baseTmpl)
}

func registerFuncs(tmpl *tmplmod.Template) {
	tmpl.Funcs(map[string]any{
		"SetUpVarName":               SetUpVarName,
		"ParseMapType":               ParseMapType,
		"ClearMapType":               ClearMapType,
		"ParseArrayType":             ParseArrayType,
		"ParseSliceType":             ParseSliceType,
		"ParsePtrType":               ParsePtrType,
		"MakeSimpleType":             MakeSimpleType,
		"MakeSimpleTypeWithAlias":    MakeSimpleTypeWithAlias,
		"MakeSimpleTypeWithEncoding": MakeSimpleTypeWithEncoding,
		"MakeValidSimpleType":        MakeValidSimpleType,
		"MakeSimpleTypeFromField":    MakeSimpleTypeFromField,
		"MakeSimplePtrTypeFromField": MakeSimplePtrTypeFromField,
		"MakeTmplData":               MakeTmplData,
		"NumSize":                    NumSize,
		"IntShift":                   IntShift,
		"ArrayIndex":                 ArrayIndex,
		"MapKeyVarName":              MapKeyVarName,
		"MapValueVarName":            MapValueVarName,
		"MakeVar":                    MakeVar,
		"MaxLastByte":                MaxLastByte,
		"include":                    MakeIncludeFunc(tmpl),
		"iterate":                    MakeIterateFunc(),
		"add":                        MakeAddFunc(),
		"minus":                      MakeMinusFunc(),
		"div":                        MakeDivFunc(),
		"replace":                    MakeReplaceFunc(),
	})
}

func registerTemplates(tmpl *tmplmod.Template) (err error) {
	var childTmpl *tmplmod.Template
	for name, template := range templates {
		childTmpl = tmpl.New(name)
		_, err = childTmpl.Parse(template)
		if err != nil {
			return
		}
	}
	return nil
}

func makeAliasFileName(lang Lang) string {
	return "alias." + lang.String() + ".tmpl"
}

func makeStructFileName(lang Lang) string {
	return "struct." + lang.String() + ".tmpl"
}
