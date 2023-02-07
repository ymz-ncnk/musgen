//go:generate go run metagen/templates_var.go
package text_template

import (
	"bytes"
	"regexp"
	tmplmod "text/template"

	"github.com/ymz-ncnk/musgen"
)

// Two main templates are templates/alias.go.tmpl and templates/struct.go.tmpl.
// After any changes in one of the template file from templates/ directory, you
// should run '$ go generate'.

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

// Generate generates code for the specified language.
func (gen MusGen) Generate(td musgen.TypeDesc, lang musgen.Lang) (
	[]byte, error) {
	err := gen.validateTypeDesc(td)
	if err != nil {
		return nil, err
	}
	switch lang {
	case musgen.GoLang:
		if musgen.Alias(td) {
			return gen.generate(td, makeAliasFileName(lang))
		}
		return gen.generate(td, makeStructFileName(lang))
	default:
		return nil, musgen.ErrUnsupportedLang
	}
}

func (gen MusGen) generate(td musgen.TypeDesc, tmplFile string) ([]byte, 
	error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := gen.baseTmpl.ExecuteTemplate(buf, tmplFile, td)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (gen MusGen) validateTypeDesc(td musgen.TypeDesc) error {
	return gen.validateEncoding(td)
}

func (gen MusGen) validateEncoding(td musgen.TypeDesc) (err error) {
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
		return musgen.ErrUnsupportedEncoding
	}
	return nil
}

func (gen MusGen) validateElemEncoding(fieldType, encoding string) error {
	var elemType string
	if at := musgen.ParseArrayType(fieldType); at.Valid {
		elemType = at.Type
	} else if st := musgen.ParseSliceType(fieldType); st.Valid {
		elemType = st.Type
	} else if mt := musgen.ParseMapType(fieldType); mt.Valid {
		elemType = mt.Value
	} else {
		return musgen.ErrUnsupportedElemEncoding
	}
	if !gen.supportEncoding(elemType, encoding) {
		return musgen.ErrUnsupportedElemEncoding
	}
	return nil
}

func (gen MusGen) validateKeyEncoding(fieldType, encoding string) error {
	var keyType string
	if mt := musgen.ParseMapType(fieldType); mt.Valid {
		keyType = mt.Key
	} else {
		return musgen.ErrUnsupportedKeyEncoding
	}
	if !gen.supportEncoding(keyType, encoding) {
		return musgen.ErrUnsupportedKeyEncoding
	}
	return nil
}

func (gen MusGen) supportEncoding(baseTmpl string, encoding string) bool {
	if encoding == musgen.RawEncoding && gen.supportRawEncoding(baseTmpl) {
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
		"SetUpVarName":               musgen.SetUpVarName,
		"ParseMapType":               musgen.ParseMapType,
		"ClearMapType":               musgen.ClearMapType,
		"ParseArrayType":             musgen.ParseArrayType,
		"ParseSliceType":             musgen.ParseSliceType,
		"ParsePtrType":               musgen.ParsePtrType,
		"MakeSimpleType":             musgen.MakeSimpleType,
		"MakeSimpleTypeWithAlias":    musgen.MakeSimpleTypeWithAlias,
		"MakeSimpleTypeWithEncoding": musgen.MakeSimpleTypeWithEncoding,
		"MakeValidSimpleType":        musgen.MakeValidSimpleType,
		"MakeSimpleTypeFromField":    musgen.MakeSimpleTypeFromField,
		"MakeSimplePtrTypeFromField": musgen.MakeSimplePtrTypeFromField,
		"MakeTmplData":               musgen.MakeTmplData,
		"NumSize":                    musgen.NumSize,
		"IntShift":                   musgen.IntShift,
		"ArrayIndex":                 musgen.ArrayIndex,
		"MapKeyVarName":              musgen.MapKeyVarName,
		"MapValueVarName":            musgen.MapValueVarName,
		"MakeVar":                    musgen.MakeVar,
		"MaxLastByte":                musgen.MaxLastByte,
		"include":                    musgen.MakeIncludeFunc(tmpl),
		"iterate":                    musgen.MakeIterateFunc(),
		"add":                        musgen.MakeAddFunc(),
		"minus":                      musgen.MakeMinusFunc(),
		"div":                        musgen.MakeDivFunc(),
		"replace":                    musgen.MakeReplaceFunc(),
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

func makeAliasFileName(lang musgen.Lang) string {
	return "alias." + lang.String() + ".tmpl"
}

func makeStructFileName(lang musgen.Lang) string {
	return "struct." + lang.String() + ".tmpl"
}
