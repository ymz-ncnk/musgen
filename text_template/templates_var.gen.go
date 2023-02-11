package text_template

var templates map[string]string

func init() {
	templates = make(map[string]string)
	templates["alias.go.tmpl"] = `// Code generated by musgen. DO NOT EDIT.

{{/* TypeDesc */}}
package {{.Package}}

// Marshal{{.Suffix}} fills buf with the MUS encoding of v.
func (v {{.Name}}) Marshal{{.Suffix}}(buf []byte) int {
  i := 0
  {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeFromField (index .Fields 0) .Unsafe .Suffix) "v") "marshal") }}
  return i
}

// Unmarshal{{.Suffix}} parses the MUS-encoded buf, and sets the result to *v.
func (v *{{.Name}}) Unmarshal{{.Suffix}}(buf []byte) (int, error) {
  i := 0
  var err error
  {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimplePtrTypeFromField (index .Fields 0) .Unsafe .Suffix) "v") "unmarshal") }}
  return i, err
}

// Size{{.Suffix}} returns the size of the MUS-encoded v.
func (v {{.Name}}) Size{{.Suffix}}() int {
  size := 0
  {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeFromField (index .Fields 0) .Unsafe .Suffix) "v") "size") }}
  return size
}`
	templates["array_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  for _, item := range {{$vn}} {
    {{- $at := (ParseArrayType $pt.Type) }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $at.Type .Unsafe .Suffix .ElemEncoding) "item") "marshal") }}
  }
}`
	templates["array_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $at := (ParseArrayType $pt.Type) }}
  for _, item := range {{$vn}} {
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $at.Type .Unsafe .Suffix .ElemEncoding) "item") "size") }}
  }
}`
	templates["array_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $at := (ParseArrayType $pt.Type)}}
  {{- $index := ArrayIndex .VarName}}
  for {{$index}} := 0; {{$index}} < {{$at.Length}}; {{$index}}++ {
    {{- $elvn := (print $vn "[" $index "]") }}
    {{- if (ParsePtrType $at.Type).Valid }}
      {{- include "initptrvar.go.tmpl" (MakeVar $elvn $at.Type false) }}
    {{- end }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeValidSimpleType $at.Type .ElemValidator 0 .Unsafe .Suffix .ElemEncoding) $elvn) "unmarshal") }}
    if err != nil {
      err = errs.NewArrayError({{$index}}, err)
      break
    }
  }
}`
	templates["bool_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  if {{$vn}} {
    buf[i] = 0x01
  } else {
    buf[i] = 0x00
  }
  i++
}`
	templates["bool_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  _ = {{.VarName}}
  size++
}`
	templates["bool_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  if i > len(buf) - 1 {
    return i, errs.ErrSmallBuf
  }
  if buf[i] == 0x01 {
    {{$vn}} = true
    i++
  } else if buf[i] == 0x00 {
    {{$vn}} = false
    i++
  } else {
    err = errs.ErrWrongByte
  }
}`
	templates["byte_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  buf[i] = byte({{$vn}})
  i++
}`
	templates["byte_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  _ = {{.VarName}}
  size++
}`
	templates["byte_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}
  {{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  if i > len(buf) - 1 {
    return i, errs.ErrSmallBuf
  }
  {{$vn}} = {{$ct}}(buf[i])
  i++
}`
	templates["custom_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  si := {{$vn}}.Marshal{{.Suffix}}(buf[i:])
  i += si
}`
	templates["custom_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  ss := {{$vn}}.Size{{.Suffix}}()
  size += ss
}`
	templates["custom_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  var sv {{$pt.Type}}
  si := 0
  si, err = sv.Unmarshal{{.Suffix}}(buf[i:])
  if err == nil {
    {{$vn}} = sv
    i += si
  }
}`
	templates["float_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $uvt := print "uint" (NumSize .Type) }}
  uv := math.Float{{(NumSize .Type)}}bits(float{{(NumSize .Type)}}({{$vn}}))
  {{- if eq .Encoding "raw" }}
    {{ include "uint_marshal_raw.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv") }}
  {{- else }}
    {{- if eq (NumSize .Type) 64 }}
      uv = (uv << 32) | (uv >> 32)
      uv = ((uv << 16) & 0xFFFF0000FFFF0000) | ((uv >> 16) & 0x0000FFFF0000FFFF)
      uv = ((uv << 8) & 0xFF00FF00FF00FF00) | ((uv >> 8) & 0x00FF00FF00FF00FF)
    {{- else if eq (NumSize .Type) 32 }}
      uv = (uv << 16) | (uv >> 16)
      uv = ((uv << 8) & 0xFF00FF00) | ((uv >> 8) & 0x00FF00FF)
    {{- end }}
    {{ include "uint_marshal.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv") }}
  {{- end }}
}`
	templates["float_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $uvt := print "uint" (NumSize .Type) }}
  {{- if eq .Encoding "raw" }}
    {{ include "uint_size_raw.go.tmpl" (SetUpVarName (MakeSimpleType  $uvt .Unsafe .Suffix) $vn) }}
  {{- else }}
    uv := math.Float{{(NumSize .Type)}}bits(float{{(NumSize .Type)}}({{$vn}}))
    {{- if eq (NumSize .Type) 64 }}
      uv = (uv << 32) | (uv >> 32)
      uv = ((uv << 16) & 0xFFFF0000FFFF0000) | ((uv >> 16) & 0x0000FFFF0000FFFF)
      uv = ((uv << 8) & 0xFF00FF00FF00FF00) | ((uv >> 8) & 0x00FF00FF00FF00FF)
    {{- else if eq (NumSize .Type) 32 }}
      uv = (uv << 16) | (uv >> 16)
      uv = ((uv << 8) & 0xFF00FF00) | ((uv >> 8) & 0x00FF00FF)
    {{- end }}
    {{ include "uint_size.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv") }}
  {{- end }}
}`
	templates["float_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}
  {{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- $uvt := print "uint" (NumSize .Type) }}
  var uv {{$uvt}}
  {{- if eq .Encoding "raw" }}
    {{ include "uint_unmarshal_raw.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv") }}
  {{- else }}
    {{ include "uint_unmarshal.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv") }}
    {{- if eq (NumSize .Type) 64 }}
      uv = (uv << 32) | (uv >> 32)
      uv = ((uv << 16) & 0xFFFF0000FFFF0000) | ((uv >> 16) & 0x0000FFFF0000FFFF)
      uv = ((uv << 8) & 0xFF00FF00FF00FF00) | ((uv >> 8) & 0x00FF00FF00FF00FF)
    {{- else if eq (NumSize .Type) 32 }}
      uv = (uv << 16) | (uv >> 16)
      uv = ((uv << 8) & 0xFF00FF00) | ((uv >> 8) & 0x00FF00FF)
    {{- end }}
  {{- end }}
  {{$vn}} = {{$ct}}(math.Float{{(NumSize .Type)}}frombits(uv))
}`
	templates["initptrvar.go.tmpl"] = `{{- /* {Name, Type, Init} */ -}}
{{- $pt := (ParsePtrType .Type) }}
{{- if $pt.Valid }}
  {{- if eq (len $pt.Stars) 1 }}
    {{.Name}} {{if .Init}}:{{end}}= new({{ClearMapType $pt.Type}})
  {{- else }}
    {{- if .Init }}
      var {{.Name}} {{ClearMapType .Type}}
    {{- end }}
    {
      tmp0 := new({{ClearMapType $pt.Type}})
      {{- range $i := (iterate (minus (len $pt.Stars) 2)) }}
        tmp{{add $i 1}} := &tmp{{$i}}
      {{- end }}
      {{.Name}} = &tmp{{minus (len $pt.Stars) 2}}
    }
  {{- end }}
{{- end }}`
	templates["int_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}{{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- if eq .Encoding "raw" }}
    {{ include "uint_marshal_raw.go.tmpl" (SetUpVarName (MakeSimpleTypeWithAlias .Alias $pt.Type .Unsafe .Suffix) $vn) }}
  {{- else }}
    {{- $uvt := print "uint" (NumSize .Type) }}
    uv := {{$uvt}}({{$vn}})
		if {{$vn}} < 0 {
			uv = ^(uv << 1)
		} else {
			uv = uv << 1
		}
    {{ include "uint_marshal.go.tmpl" (SetUpVarName (MakeSimpleTypeWithAlias .Alias $pt.Type .Unsafe .Suffix) "uv") }}
  {{- end }}
}`
	templates["int_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}{{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- if eq .Encoding "raw" }}
    {{ include "uint_size_raw.go.tmpl" (SetUpVarName (MakeSimpleType  $pt.Type .Unsafe .Suffix) $vn) }}
  {{- else }}
    {{- $uvt := print "uint" (NumSize .Type) }}
    uv := {{$uvt}}({{$vn}}<<1) ^ {{$uvt}}({{$vn}}>>{{IntShift .Type}})
    {{ include "uint_size.go.tmpl" (SetUpVarName (MakeSimpleTypeWithAlias .Alias $pt.Type .Unsafe .Suffix) "uv") }}
  {{- end }}
}`
	templates["int_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}{{if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- if eq .Encoding "raw" }}
    {{ include "uint_unmarshal_raw.go.tmpl" (SetUpVarName (MakeSimpleTypeWithAlias .Alias $pt.Type .Unsafe .Suffix) $vn) }}
  {{- else }}
    {{- $uvt := print "uint" (NumSize .Type) }}
    var uv {{$uvt}}
    {{ include "uint_unmarshal.go.tmpl" (SetUpVarName (MakeSimpleType $uvt .Unsafe .Suffix) "uv")}}
		if uv&1 == 1 {
			uv = ^(uv >> 1)
		} else {
			uv = uv >> 1
		}
    {{$vn}} = {{$ct}}(uv)
  {{- end }}
}`
	templates["map_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_marshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  for ke, vl := range {{$vn}} {
    {{- $mt := (ParseMapType $pt.Type) }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $mt.Key .Unsafe .Suffix .KeyEncoding) "ke") "marshal") }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $mt.Value .Unsafe .Suffix .ElemEncoding) "vl") "marshal") }}
  }
}`
	templates["map_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_size.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  {{- $mt := (ParseMapType $pt.Type) }}
  for ke, vl := range {{$vn}} {
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $mt.Key .Unsafe .Suffix .KeyEncoding) "ke") "size") }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $mt.Value .Unsafe .Suffix .ElemEncoding) "vl") "size") }}
  }
}`
	templates["map_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  var length int
  {{ include "int_unmarshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  if length < 0 {
    return i, errs.ErrNegativeLength
  }
  {{- if ne .MaxLength 0 }}
    if length > {{.MaxLength}} {
      err = errs.ErrMaxLengthExceeded
    } else {
  {{- end }}
    {{$vn}} = make({{ClearMapType $pt.Type}})
    {{- $mpt := (ParseMapType $pt.Type)}}
    {{- $kevn := MapKeyVarName .VarName }}
    {{- $vlvn := MapValueVarName .VarName }}
    for ; length > 0; length-- {
      {{- if (ParsePtrType $mpt.Key).Valid }}
        {{- include "initptrvar.go.tmpl" (MakeVar $kevn $mpt.Key true) }}
      {{- else }}
        var {{$kevn}} {{ClearMapType $mpt.Key}}
      {{- end }}
      {{- if (ParsePtrType $mpt.Value).Valid }}
        {{- include "initptrvar.go.tmpl" (MakeVar $vlvn $mpt.Value true) }}
      {{- else }}
        var {{$vlvn}} {{ClearMapType $mpt.Value}}
      {{- end }}
      {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeValidSimpleType $mpt.Key .KeyValidator 0 .Unsafe .Suffix .KeyEncoding) $kevn) "unmarshal") }}
      if err != nil {
        err = errs.NewMapKeyError({{$kevn}}, err)
        break
      }
      {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeValidSimpleType $mpt.Value .ElemValidator 0 .Unsafe .Suffix .ElemEncoding) $vlvn) "unmarshal") }}
      if err != nil {
        err = errs.NewMapValueError({{$kevn}}, {{$vlvn}}, err)
        break
      }
      ({{$vn}})[{{$kevn}}] = {{$vlvn}}
    }
  {{- if ne .MaxLength 0 }}
    }
  {{- end }}
}`
	templates["ptr_end.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
}`
	templates["ptr_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{{- $vn := .VarName -}}
if {{$vn}} == nil {
	buf[i] = 0
	i++
} else {
	buf[i] = 1
	i++`
	templates["ptr_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{{- $vn := .VarName -}}
size++
if {{$vn}} != nil {`
	templates["ptr_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{{- $vn := .VarName -}}
if buf[i] == 0 {
	i++
	{{$vn}} = nil
} else if buf[i] != 1 {
	i++
	return i, errs.ErrWrongByte
} else {
	i++`
	templates["simpletypes.go.tmpl"] = `{{- /* {SimpleTypeVar, MUName} */ -}}
{{- $pt := (ParsePtrType .SimpleTypeVar.Type)}}
{{- if $pt.Valid }}
  {{- if eq .MUName "marshal" }}
    {{- include "ptr_marshal.go.tmpl" . }}
  {{ else if and (eq .MUName "unmarshal") (eq .SimpleTypeVar.Alias "") }}
    {{- include "ptr_unmarshal.go.tmpl" . }}
  {{ else if eq .MUName "size" }}
    {{- include "ptr_size.go.tmpl" . }}
  {{ end }}
{{- end }}
  {{- if eq $pt.Type "uint64" "uint32" "uint16" "uint" }}
    {{- if eq .SimpleTypeVar.Encoding "raw" }}
      {{- include (print "uint_" .MUName "_raw.go.tmpl") .SimpleTypeVar }}
    {{- else }}
      {{- include (print "uint_" .MUName ".go.tmpl") .SimpleTypeVar }}
    {{- end }}
  {{- else if eq $pt.Type "int64" "int32" "int16" "int" }}
    {{- include (print "int_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if eq $pt.Type "float64" "float32" }}
    {{- include (print "float_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if eq $pt.Type "string" }}
    {{- include (print "string_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if eq $pt.Type "bool" }}
    {{- include (print "bool_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if eq $pt.Type "byte" "uint8" "int8" }}
    {{- include (print "byte_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if (ParseMapType .SimpleTypeVar.Type).Valid }}
    {{- include (print "map_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if (ParseArrayType .SimpleTypeVar.Type).Valid }}
    {{- include (print "array_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else if (ParseSliceType .SimpleTypeVar.Type).Valid }}
    {{- include (print "slice_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- else }}
    {{- include (print "custom_" .MUName ".go.tmpl") .SimpleTypeVar }}
  {{- end }}
{{- if $pt.Valid }}
  {{- if or (eq .MUName "marshal") (and (eq .MUName "unmarshal") (eq .SimpleTypeVar.Alias "")) (eq .MUName "size") }}
    {{- include "ptr_end.go.tmpl" . }}
  {{- end }}
{{- end }}
{{- if eq .MUName "unmarshal" }}
  {{- include "validator.go.tmpl" . }}
{{- end }}`
	templates["slice_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_marshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  for _, el := range {{$vn}} {
    {{- $st := ParseSliceType $pt.Type }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $st.Type .Unsafe .Suffix .ElemEncoding) "el") "marshal") }}
  }
}`
	templates["slice_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_size.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  for _, el := range {{$vn}} {
    {{- $st := ParseSliceType $pt.Type }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeWithEncoding $st.Type .Unsafe .Suffix .ElemEncoding) "el") "size") }}
  }
}`
	templates["slice_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $st := (ParseSliceType $pt.Type).Type}}
  {{- $index := ArrayIndex .VarName}}
  var length int
  {{ include "int_unmarshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  if length < 0 {
    return i, errs.ErrNegativeLength
  }
  {{- if ne .MaxLength 0 }}
    if length > {{.MaxLength}} {
      err = errs.ErrMaxLengthExceeded
    } else {
  {{- end }}
    {{$vn}} = make({{ClearMapType $pt.Type}}, length)
    for {{$index}} := 0; {{$index}} < length; {{$index}}++ {
      {{- $elvn := (print $vn "[" $index "]") }}
      {{- if (ParsePtrType $st).Valid }}
        {{- include "initptrvar.go.tmpl" (MakeVar $elvn $st false) }}
      {{- end }}
      {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeValidSimpleType $st .ElemValidator 0 .Unsafe .Suffix .ElemEncoding) (print $vn "[" $index "]")) "unmarshal") }}
      if err != nil {
        err = errs.NewSliceError({{$index}}, err)
        break
      }
    }
  {{- if ne .MaxLength 0 }}
    }
  {{- end }}
}`
	templates["string_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_marshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  if len(buf[i:]) < length {
    panic(errs.ErrSmallBuf)
  }
  i += copy(buf[i:], {{$vn}})
}`
	templates["string_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  length := len({{$vn}})
  {{ include "int_size.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  size += len({{$vn}})
}`
	templates["string_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type) }}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type }}
  {{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  var length int
  {{ include "int_unmarshal.go.tmpl" (SetUpVarName (MakeSimpleType "int" .Unsafe .Suffix) "length") }}
  if length < 0 {
    return i, errs.ErrNegativeLength
  } 
  if len(buf) < i+length {
    return i, errs.ErrSmallBuf
  }
  {{- if ne .MaxLength 0 }}
    if length > {{.MaxLength}} {
      err = errs.ErrMaxLengthExceeded
    } else {
  {{- end }}
    {{- if .Unsafe }}
      content := buf[i : i+length]
      slcHeader := (*reflect.SliceHeader)(unsafe.Pointer(&content))
      {{- if eq (len $pt.Stars) 0 }}
        strHeader := (*reflect.StringHeader)(unsafe.Pointer(&{{.VarName}}))
      {{- else if eq (len $pt.Stars) 1 }}
        strHeader := (*reflect.StringHeader)(unsafe.Pointer({{.VarName}}))
      {{- else }}
        strHeader := (*reflect.StringHeader)(unsafe.Pointer({{replace $pt.Stars "*" "" 1}}{{.VarName}}))
      {{- end }}
      strHeader.Data = slcHeader.Data
      strHeader.Len = slcHeader.Len
    {{- else }}
      {{$vn}} = {{$ct}}(buf[i : i+length])
    {{- end}}
    i += length
  {{- if ne .MaxLength 0 }}
    }
  {{- end }}
}`
	templates["struct.go.tmpl"] = `// Code generated by musgen. DO NOT EDIT.

{{/* TypeDesc */}}
package {{.Package}}

// Marshal{{.Suffix}} fills buf with the MUS encoding of v.
func (v {{.Name}}) Marshal{{.Suffix}}(buf []byte) int {
  i := 0
  {{- $unsafe := .Unsafe }}
  {{- $suffix := .Suffix}}
  {{- range $index, $field := .Fields }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeFromField $field $unsafe $suffix) (print "v." $field.Name)) "marshal") }}
  {{- end }}
  return i
}

// Unmarshal{{.Suffix}} parses the MUS-encoded buf, and sets the result to *v.
func (v *{{.Name}}) Unmarshal{{.Suffix}}(buf []byte) (int, error) {
  i := 0
  {{- $unsafe := .Unsafe }}
  {{- $suffix := .Suffix}}
  var err error
  {{- range $index, $field := .Fields }}
    {{- $fvn := (print "v." $field.Name) }}
    {{- if (ParsePtrType $field.Type).Valid }}
      {{- include "initptrvar.go.tmpl" (MakeVar $fvn $field.Type false) }}
    {{- end }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeFromField $field $unsafe $suffix) $fvn) "unmarshal") }}
    if err != nil {
      return i, errs.NewFieldError({{print "\"" $field.Name "\""}}, err)
    }
  {{- end }}
  return i, err
}

// Size{{.Suffix}} returns the size of the MUS-encoded v.
func (v {{.Name}}) Size{{.Suffix}}() int {
  size := 0
  {{- $unsafe := .Unsafe }}
  {{- $suffix := .Suffix}}
  {{- range $index, $field := .Fields }}
    {{ include "simpletypes.go.tmpl" (MakeTmplData (SetUpVarName (MakeSimpleTypeFromField $field $unsafe $suffix) (print "v." $field.Name)) "size") }}
  {{- end }}
  return size
}`
	templates["uint_marshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}
		tmp := {{print "(" $pt.Stars .VarName ")"}}
		{{- $vn = "tmp"}}
	{{- end }}	
  for {{$vn}} >= 0x80 {
    buf[i] = byte({{$vn}}) | 0x80
    {{$vn}} >>= 7
    i++
  }
  buf[i] = byte({{$vn}})
  i++
}`
	templates["uint_marshal_raw.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- $ct := $pt.Type}}
  {{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- $ns := NumSize $pt.Type}}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- if .Unsafe }}
    *(*{{$ct}})(unsafe.Pointer(&buf[i])) = {{$vn}}
    i += {{div $ns 8}}
  {{- else }}
    buf[i] = byte({{$vn}})
    i++
    {{- if ge $ns 16 }}
      buf[i] = byte({{$vn}} >> 8)
      i++
    {{- end }}      
    {{- if ge $ns 32 }}
      buf[i] = byte({{$vn}} >> 16)
      i++
      buf[i] = byte({{$vn}} >> 24)
      i++
    {{- end }}  
    {{- if eq $ns 64 }}
      buf[i] = byte({{$vn}} >> 32)
      i++
      buf[i] = byte({{$vn}} >> 40)
      i++
      buf[i] = byte({{$vn}} >> 48)
      i++
      buf[i] = byte({{$vn}} >> 56)
      i++
    {{- end }}
  {{- end }}
}`
	templates["uint_size.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}
		tmp := {{print "(" $pt.Stars .VarName ")"}}
		{{- $vn = "tmp"}}
	{{- end }}
  for {{$vn}} >= 0x80 {
    {{$vn}} >>= 7
    size++
  }
  size++
}`
	templates["uint_size_raw.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $ns := NumSize $pt.Type}}
  _ = {{.VarName}}
  size += {{div $ns 8}}
}`
	templates["uint_unmarshal.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}{{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  if i > len(buf) - 1 {
    return i, errs.ErrSmallBuf
  }
  shift := 0
  done := false
  for l, b := range buf[i:] {
    if l == {{minus .MaxLength 1}} && b > {{MaxLastByte $pt.Type}} {
      return i, errs.ErrOverflow
    }
    if b < 0x80 {
      {{$vn}} = {{$vn}} | {{$ct}}(b)<<shift
      done = true
      i += l+1
      break
    }
    {{$vn}} = {{$vn}} | {{$ct}}(b&0x7F)<<shift
    shift += 7
  }
  if !done {
    return i, errs.ErrSmallBuf
  }
}`
	templates["uint_unmarshal_raw.go.tmpl"] = `{{- /* SimpleTypeVar */ -}}
{
  {{- $pt := (ParsePtrType .Type)}}
  {{- $vn := .VarName }}
  {{- if $pt.Valid }}{{ $vn = print "(" $pt.Stars .VarName ")" }}{{ end }}
  {{- $ct := $pt.Type}}{{- if ne .Alias "" }}{{$ct = .Alias}}{{ end }}
  {{- $ns := NumSize $pt.Type}}
  if len(buf) < {{div $ns 8}} {
    return i, errs.ErrSmallBuf
  }
  {{- if .Unsafe }}
    {{$vn}} = *(*{{$ct}})(unsafe.Pointer(&buf[i]))
    i += {{div $ns 8}}
  {{- else }}
    {{$vn}} = {{$ct}}(buf[i])
    i++
    {{- if ge $ns 16 }}
      {{$vn}} |= {{$ct}}(buf[i])<<8
      i++
    {{- end }}
    {{- if ge $ns 32 }}
      {{$vn}} |= {{$ct}}(buf[i])<<16
      i++
      {{$vn}} |= {{$ct}}(buf[i])<<24
      i++
    {{- end }}
    {{- if eq $ns 64 }}
      {{$vn}} |= {{$ct}}(buf[i])<<32
      i++
      {{$vn}} |= {{$ct}}(buf[i])<<40
      i++
      {{$vn}} |= {{$ct}}(buf[i])<<48
      i++
      {{$vn}} |= {{$ct}}(buf[i])<<56
      i++
    {{- end }}
  {{- end }}
}`
	templates["validator.go.tmpl"] = `{{- if ne .Validator "" }}
  if err == nil {
    err = {{.Validator}}({{.VarName}})
  }
{{- end }}`
}
