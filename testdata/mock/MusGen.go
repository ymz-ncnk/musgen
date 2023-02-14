package mocks

import (
	"github.com/ymz-ncnk/amock"
	"github.com/ymz-ncnk/musgen/v2"
)

func NewMusGen() MusGen {
	return MusGen{amock.New("MusGen")}
}

type MusGen struct {
	*amock.Mock
}

func (musGen MusGen) RegisterGenerate(fn func(td musgen.TypeDesc,
	lang musgen.Lang) ([]byte, error)) MusGen {
	musGen.Register("Generate", fn)
	return musGen
}

func (musGen MusGen) Generate(td musgen.TypeDesc, lang musgen.Lang) (
	data []byte, err error) {
	vals, err := musGen.Call("Generate", td, lang)
	if err != nil {
		return
	}
	data, _ = vals[0].([]byte)
	err, _ = vals[1].(error)
	return
}
