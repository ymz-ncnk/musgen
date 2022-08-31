// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	templatesDir      = "templates/"
	tmplsInitFileName = "tmpls.gen.go"
	fdescFileName     = "fdesc.musgen.go"
	tdescFileName     = "tdesc.musgen.go"
	tmplsInit         = "" +
		"package musgen\n\n" +
		"var tmpls map[string]string\n" +
		"func init() {\n" +
		"	tmpls = make(map[string]string)\n" +
		"	{{- range $tmplName, $tmpl := . }}\n" +
		"		tmpls[\"{{$tmplName}}\"] = `{{$tmpl}}`\n" +
		"	{{- end }}\n" +
		"}"
)

// Transforms files from templates/ dir to map[string]string
// representation, where a key is a filename, a value is a file's content.
// Generates tmplsinit.gen.go.
func main() {
	err := genTmpls()
	if err != nil {
		panic(err)
	}
}

func genTmpls() error {
	data, err := makeTmplsInitData()
	if err != nil {
		return err
	}
	var bs []byte
	bs, err = makeTmplsInitFile(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tmplsInitFileName, bs, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func makeTmplsInitData() (map[string]string, error) {
	data := make(map[string]string)
	dirs, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}
	var fs []os.FileInfo
	var content []byte
	for i := 0; i < len(dirs); i++ {
		if !dirs[i].IsDir() {
			return nil, fmt.Errorf("found not dir file %v", dirs[i].Name())
		}
		fs, err = ioutil.ReadDir(templatesDir + dirs[i].Name())
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(fs); j++ {
			content, err = ioutil.ReadFile(templatesDir + dirs[i].Name() + "/" +
				fs[j].Name())
			if err != nil {
				return nil, err
			}
			data[string(fs[j].Name())] = string(content)
		}
	}
	return data, nil
}

func makeTmplsInitFile(data map[string]string) ([]byte, error) {
	t := template.New("base")
	t.Parse(tmplsInit)
	buf := bytes.NewBuffer(make([]byte, 0))
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
