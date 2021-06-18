package musgen

import (
	"testing"
)

func TestGoLangAliasMusgen(t *testing.T) {
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

func TestGoLangStructMusgen(t *testing.T) {
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
