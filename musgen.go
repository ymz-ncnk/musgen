//go:generate go run metagen/metagen.go
package musgen

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"text/template"
)

// Two main templates are alias.go.tmpl and struct.go.tmpl.
// After changing a template run '$ go generate'.

// UintWithSystemSize system uint type with size.
var UintWithSystemSize = "uint" + strconv.Itoa(strconv.IntSize)

const (
	MaxZigZagLength64 = 10
	MaxZigZagLength32 = 5
	MaxZigZagLength16 = 3
	MaxZigZagLength8  = 1
	RawEncoding       = "raw"
)

// New returns a new MusGen.
func New() (MusGen, error) {
	t := template.New("base")
	funcs := make(template.FuncMap)
	funcs["SetUpVarName"] = SetUpVarName
	funcs["ParseMapType"] = ParseMapType
	funcs["ClearMapType"] = ClearMapType
	funcs["ParseArrayType"] = ParseArrayType
	funcs["ParseSliceType"] = ParseSliceType
	funcs["ParsePtrType"] = ParsePtrType
	funcs["MakeSimpleType"] = MakeSimpleType
	funcs["MakeSimpleTypeWithAlias"] = MakeSimpleTypeWithAlias
	funcs["MakeSimpleTypeWithEncoding"] = MakeSimpleTypeWithEncoding
	funcs["MakeValidSimpleType"] = MakeValidSimpleType
	funcs["MakeSimpleTypeFromField"] = MakeSimpleTypeFromField
	funcs["MakeSimplePtrTypeFromField"] = MakeSimplePtrTypeFromField
	funcs["MakeTmplData"] = MakeTmplData
	funcs["NumSize"] = NumSize
	funcs["IntShift"] = IntShift
	funcs["ArrayIndex"] = ArrayIndex
	funcs["MapKeyVarName"] = MapKeyVarName
	funcs["MapValueVarName"] = MapValueVarName
	funcs["MakeVar"] = MakeVar
	funcs["MaxLastByte"] = MaxLastByte
	funcs["include"] = IncludeFunc(t)
	funcs["iterate"] = IterateFunc()
	funcs["add"] = AddFunc()
	funcs["minus"] = MinusFunc()
	funcs["div"] = DivFunc()
	funcs["replace"] = ReplaceFunc()
	t.Funcs(funcs)
	var tl *template.Template
	var err error
	for tmplName, tmpl := range tmpls {
		tl = t.New(tmplName)
		_, err = tl.Parse(tmpl)
		if err != nil {
			return MusGen{}, err
		}
	}
	return MusGen{t}, err
}

// MusGen is a code generator for the MUS format.
type MusGen struct {
	t *template.Template
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
			return gen.generate(td, "alias."+lang.String()+".tmpl")
		}
		return gen.generate(td, "struct."+lang.String()+".tmpl")
	default:
		return nil, fmt.Errorf("unsuported lang %v", lang.String())
	}
}

func (gen MusGen) generate(td TypeDesc, tmplFile string) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := gen.t.ExecuteTemplate(buf, tmplFile, td)
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

func (gen MusGen) supportEncoding(t string, encoding string) bool {
	if encoding == RawEncoding && gen.supportRawEncoding(t) {
		return true
	}
	return false
}

func (gen MusGen) supportRawEncoding(t string) bool {
	re := regexp.MustCompile(`^\**(?:(?:uint|int)(?:64|32|16|8|)|float(?:64|32))$`)
	return re.MatchString(t)
}
