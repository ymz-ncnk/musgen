//go:generate go run metagen/metagen.go
package musgen

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// Two main templates are alias.go.tmpl and struct.go.tmpl.
// After changing a template run '$ go generate'.

// UintWithSystemSize system uint type with size.
var UintWithSystemSize = "uint" + strconv.Itoa(strconv.IntSize)

// Lang represents supported languages.
type Lang int

const (
	// GoLang constant.
	GoLang Lang = iota
)

func (l Lang) String() string {
	return [...]string{"go"}[l]
}

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
	KeyValidator  string
	ElemValidator string
}

// SimpleTypeVar pipeline for the simpletypes.go.tmpl.
type SimpleTypeVar struct {
	SimpleType
	VarName string
}

// SimpleType represents a simplest supported type, like bool, uint64, ...
type SimpleType struct {
	Type          string
	Suffix        string
	Unsafe        bool
	MaxLength     int
	Alias         string
	Validator     string
	KeyValidator  string
	ElemValidator string
}

var tmpls map[string]string

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

// Alias checks if TypeDesc is an Alias type.
func Alias(td TypeDesc) bool {
	return len(td.Fields) == 1 && td.Fields[0].Name == "" &&
		td.Name == td.Fields[0].Alias
}

// MakeTmplData is a template function. Creates a pipeline for the
// simpletypes.go.tmpl.
func MakeTmplData(simpleTypeVar SimpleTypeVar, muName string) struct {
	SimpleTypeVar
	MUName string
} {
	return struct {
		SimpleTypeVar
		MUName string
	}{
		SimpleTypeVar: simpleTypeVar,
		MUName:        muName,
	}
}

// SetUpVarName is a template function. Sets up variable name for the
// SimpleType.
func SetUpVarName(simpleType SimpleType, varName string) SimpleTypeVar {
	return SimpleTypeVar{
		SimpleType: simpleType,
		VarName:    varName,
	}
}

// ParseMapType is a template function. If the given string represents a map
// type, Valid == true. The required format is map-n[Key]Value-n, where n is a
// map number.
//
// Map number helps to parse a map type. For example, we could rewrite the type
// map[string]map[string]int as - map-1[string]-1map-0[string]-0int.
func ParseMapType(t string) struct {
	Valid bool
	Key   string
	Value string
} {
	pre := regexp.MustCompile(`^\**map(-\d+)`)
	match := pre.FindStringSubmatch(t)
	if len(match) == 2 {
		mapNum := match[1]
		re := regexp.MustCompile(`^\**map` + mapNum + `\[(.+?)\]` + mapNum + `(.+$)`)
		match := re.FindStringSubmatch(t)
		if len(match) == 3 {
			return struct {
				Valid bool
				Key   string
				Value string
			}{
				Valid: true,
				Key:   match[1],
				Value: match[2],
			}
		}
	}
	return struct {
		Valid bool
		Key   string
		Value string
	}{
		Valid: false,
		Key:   "",
		Value: "",
	}
}

// ClearMapType is a template function. It removes map numbers.
func ClearMapType(mt string) string {
	re := regexp.MustCompile(`-\d`)
	return re.ReplaceAllString(mt, "")
}

// ParseArrayType is a template function. If the given string represents an
// array type, than Valid == true. The required format is [Length]Type.
func ParseArrayType(t string) struct {
	Valid  bool
	Type   string
	Length int
} {
	re := regexp.MustCompile(`^\**\[(\d+)\](.+$)`)
	match := re.FindStringSubmatch(t)
	if len(match) == 3 {
		length, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		return struct {
			Valid  bool
			Type   string
			Length int
		}{
			Valid:  true,
			Type:   match[2],
			Length: length,
		}
	}
	return struct {
		Valid  bool
		Type   string
		Length int
	}{
		Valid:  false,
		Type:   "",
		Length: 0,
	}
}

// ParseSliceType is a template function. If the given string represents a
// slice type, than Valid == true. The required format is []Type.
func ParseSliceType(t string) struct {
	Valid bool
	Type  string
} {
	re := regexp.MustCompile(`^\**\[\](.+$)`)
	match := re.FindStringSubmatch(t)
	if len(match) == 2 {
		return struct {
			Valid bool
			Type  string
		}{
			Valid: true,
			Type:  match[1],
		}
	}
	return struct {
		Valid bool
		Type  string
	}{
		Valid: false,
		Type:  "",
	}
}

// ParsePtrType is a template function. If the given string represents a
// pointer type, than Valid == true. The required format is *Type.
//
// If Valid == false, Type is equal to the given type.
func ParsePtrType(t string) struct {
	Valid bool
	Stars string
	Type  string
} {
	re := regexp.MustCompile(`(^\*+)(.+$)`)
	match := re.FindStringSubmatch(t)
	if len(match) == 3 {
		return struct {
			Valid bool
			Stars string
			Type  string
		}{
			Valid: true,
			Stars: match[1],
			Type:  match[2],
		}
	}
	return struct {
		Valid bool
		Stars string
		Type  string
	}{
		Valid: false,
		Stars: "",
		Type:  t,
	}
}

// func ParsePtrType(t string) struct {
// 	Valid bool
// 	Type  string
// } {
// 	re := regexp.MustCompile(`^\*(.+$)`)
// 	match := re.FindStringSubmatch(t)
// 	if len(match) == 2 {
// 		return struct {
// 			Valid bool
// 			Type  string
// 		}{
// 			Valid: true,
// 			Type:  match[1],
// 		}
// 	}
// 	return struct {
// 		Valid bool
// 		Type  string
// 	}{
// 		Valid: false,
// 		Type:  t,
// 	}
// }

// MakeSimpleType creates SimpleType.
func MakeSimpleType(t string, unsafe bool, suffix string) SimpleType {
	return SimpleType{
		Type:      t,
		Suffix:    suffix,
		Unsafe:    unsafe,
		MaxLength: defMaxLength(t),
	}
}

// MakeValidSimpleType creates SimpleType with the validator and maxLength.
func MakeValidSimpleType(t, validator string,
	maxLength int, unsafe bool, suffix string) SimpleType {
	st := SimpleType{
		Type:      t,
		Suffix:    suffix,
		Unsafe:    unsafe,
		Validator: validator,
		MaxLength: defMaxLength(t),
	}
	if st.MaxLength == 0 {
		st.MaxLength = maxLength
	}
	return st
}

// MakeSimpleTypeFromField creates SimpleType from FieldDesc.
func MakeSimpleTypeFromField(f FieldDesc, unsafe bool, suffix string) SimpleType {
	st := SimpleType{
		Type:          f.Type,
		Suffix:        suffix,
		Unsafe:        unsafe,
		Alias:         f.Alias,
		Validator:     f.Validator,
		KeyValidator:  f.KeyValidator,
		ElemValidator: f.ElemValidator,
		MaxLength:     defMaxLength(f.Type),
	}
	if st.MaxLength == 0 {
		st.MaxLength = f.MaxLength
	}
	return st
}

// MakeSimplePtrTypeFromField creates SimpleType from FieldDesc.
func MakeSimplePtrTypeFromField(f FieldDesc, unsafe bool,
	suffix string) SimpleType {
	f.Type = "*" + f.Type
	return MakeSimpleTypeFromField(f, unsafe, suffix)
}

// IncludeFunc creates template's include func.
func IncludeFunc(tmpl *template.Template) func(string, interface{}) (string,
	error) {
	return func(name string, pipeline interface{}) (string, error) {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, name, pipeline); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

// IterateFunc creates template's iterate func.
func IterateFunc() func(int) []int {
	return func(count int) []int {
		var items []int
		for i := 0; i < (count); i++ {
			items = append(items, i)
		}
		return items
	}
}

// AddFunc creates template's add func.
func AddFunc() func(int, int) int {
	return func(a int, b int) int {
		return a + b
	}
}

// MinusFunc creates template's minus func.
func MinusFunc() func(int, int) int {
	return func(a int, b int) int {
		return a - b
	}
}

// ReplaceFunc creates template's replace func.
func ReplaceFunc() func(string, string, string, int) string {
	return func(s string, old string, new string, n int) string {
		return strings.Replace(s, old, new, n)
	}
}

// NumSize is a template function. Returns a size of the number type.
// For int32 returns 32, for uint8 returns 8, and so on, for int returns the
// system size.
func NumSize(t string) int {
	re := regexp.MustCompile(`(\d\d)$`)
	match := re.FindStringSubmatch(t)
	if len(match) == 2 {
		n, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		return n
	}
	return strconv.IntSize
}

// IntShift is a template function. The given string should represent sized
// type, like int8, int16, ... Returns the size redused by one.
func IntShift(t string) string {
	return strconv.Itoa(NumSize(t) - 1)
}

// ArrayIndex is a template function. It makes from vn a variable name for the
// array index, which is used in a 'for' cycle.
// The first array index is j, subarray index is jj, next subarray index is
// jjj, ...
func ArrayIndex(vn string) string {
	return "j" + prevArrayIndex(vn)
}

func prevArrayIndex(vn string) string {
	re := regexp.MustCompile(`\[(j+)\]$`)
	match := re.FindStringSubmatch(vn)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

// MapKeyVarName is a template function. It makes from vn a variable name for a
// map key.
// The first map key variable name is kem, submap key variable name is kemm,
// next key variable name is kemmm, ...
func MapKeyVarName(vn string) string {
	return mapUnitVarName(vn, "ke")
}

// MapValueVarName is a template function. It makes from vn a variable name for
// a map value.
// The first map value variable name is vlm, submap value variable name is vlmm,
// next submap value variable name is vlmmm, ...
func MapValueVarName(vn string) string {
	return mapUnitVarName(vn, "vl")
}

// MakeVar returns a pipeline for the initptrvar.go.tmpl.
func MakeVar(name string, t string, init bool) struct {
	Name string
	Type string
	Init bool
} {
	return struct {
		Name string
		Type string
		Init bool
	}{
		Name: name,
		Type: t,
		Init: init,
	}
}

func mapUnitVarName(vn string, unit string) string {
	re := regexp.MustCompile("m+$")
	return unit + re.FindString(filterNotAlphabeticChars(vn)+"m")
}

func filterNotAlphabeticChars(vn string) string {
	re := regexp.MustCompile(`[^a-zA-Z]`)
	return re.ReplaceAllString(vn, "")
}

func defMaxLength(t string) int {
	re := regexp.MustCompile(`^\**`)
	t = re.ReplaceAllString(t, "")
	switch t {
	case "uint64":
		return 10
	case "uint32":
		return 5
	case "uint16":
		return 3
	case "uint8":
		return 1
	case "uint":
		return defMaxLength(UintWithSystemSize)
	}
	return 0
}

// MaxLastByte returns how big could be the last allowed byte for the
// specified type.
// For example, "uint64" number couldn't take more than 9 bytes, and the last
// one must be less or equal to 1.
func MaxLastByte(t string) int {
	switch t {
	case "uint64":
		return 1
	case "uint32":
		return 15
	case "uint16":
		return 3
	case "uint":
		return MaxLastByte(UintWithSystemSize)
	default:
		panic("unsupported type")
	}
}
