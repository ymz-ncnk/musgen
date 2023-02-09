package mocks

import (
	"github.com/ymz-ncnk/amock"
	"github.com/ymz-ncnk/musgen"
)

func NewMusGenMock() MusGenMock{
	return MusGenMock{amock.New("MusGen")}
}

type MusGenMock struct {
	*amock.Mock
}

func (musGen MusGenMock) RegisterGenerate(fn func(td musgen.TypeDesc, 
	lang musgen.Lang) ([]byte, error)) MusGenMock {
	musGen.Register("Generate", fn)
	return musGen
}

func (musGen MusGenMock) Generate(td musgen.TypeDesc, lang musgen.Lang) (
	data []byte, err error) {
	vals, err := musGen.Call("Generate", td, lang)
	if err != nil {
		return
	}
	data, _ = vals[0].([]byte)
	err, _ = vals[1].(error)
	return
}