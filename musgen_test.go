package musgen

import (
	"testing"
)

func TestGoLangAliasMusGen(t *testing.T) {
	musGen, err := New()
	if err != nil {
		t.Error(err)
	}
	td := TypeDesc{
		Package: "musgen",
		Name:    "IntAlias",
		Fields:  []FieldDesc{{Type: "int", Alias: "IntAlias"}},
		Suffix:  "MUS",
	}
	_, err = musGen.Generate(td, GoLang)
	if err != nil {
		t.Error(err)
	}
}

func TestGoLangStructMusGen(t *testing.T) {
	musGen, err := New()
	if err != nil {
		t.Error(err)
	}
	td := TypeDesc{
		Package: "musgen",
		Name:    "IntStruct",
		Fields:  []FieldDesc{{Name: "Int", Type: "int"}},
		Suffix:  "MUS",
	}
	_, err = musGen.Generate(td, GoLang)
	if err != nil {
		t.Error(err)
	}
}

func TestGoLangRawEncoding(t *testing.T) {
	musGen, err := New()
	if err != nil {
		t.Error(err)
	}
	td := TypeDesc{
		Package: "musgen",
		Name:    "Uint64RawAlias",
		Fields: []FieldDesc{{Type: "uint64", Alias: "Uint64RawAlias",
			Encoding: "raw"}},
		Unsafe: true,
		Suffix: "MUS",
	}
	_, err = musGen.Generate(td, GoLang)
	if err != nil {
		t.Error(err)
	}
}

func TestGoLangNotValidRawEncoding(t *testing.T) {
	musGen, err := New()
	if err != nil {
		t.Error(err)
	}
	{
		td := TypeDesc{
			Package: "musgen",
			Name:    "StringRawAlias",
			Fields: []FieldDesc{{Type: "string", Alias: "StringRawAlias",
				Encoding: "raw"}},
			Unsafe: true,
			Suffix: "MUS",
		}
		_, err = musGen.Generate(td, GoLang)
		if err == nil {
			t.Error("not valid encoding ok")
		}
		if err != ErrUnsupportedEncoding {
			t.Error("unexpected error")
		}
	}
	{
		td := TypeDesc{
			Package: "musgen",
			Name:    "ArrayRawAlias",
			Fields: []FieldDesc{{Type: "[2]bool", Alias: "ArrayRawAlias",
				ElemEncoding: "raw"}},
			Unsafe: true,
			Suffix: "MUS",
		}
		_, err = musGen.Generate(td, GoLang)
		if err == nil {
			t.Error("not valid encoding ok")
		}
		if err != ErrUnsupportedElemEncoding {
			t.Error("unexpected error")
		}
	}
	{
		td := TypeDesc{
			Package: "musgen",
			Name:    "MapRawAlias",
			Fields: []FieldDesc{{Type: "map-0[string]-0uint", Alias: "MapRawAlias",
				KeyEncoding: "raw", ElemEncoding: "raw"}},
			Unsafe: true,
			Suffix: "MUS",
		}
		_, err = musGen.Generate(td, GoLang)
		if err == nil {
			t.Error("not valid encoding ok")
		}
		if err != ErrUnsupportedKeyEncoding {
			t.Error("unexpected error")
		}
	}
}
