//go:build ignore

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
	templatesDir           = "templates/"
	templatesVarFileName   = "templates_var.gen.go"
	templatesVarFileSample = "" +
		"package text_template\n\n" +
		"var templates map[string]string\n" +
		"func init() {\n" +
		"	templates = make(map[string]string)\n" +
		"	{{- range $tmplName, $tmpl := . }}\n" +
		"		templates[\"{{$tmplName}}\"] = `{{$tmpl}}`\n" +
		"	{{- end }}\n" +
		"}"
)

// Transforms files from templates/ dir to map[string]string
// representation, where a key is a filename, a value is a file's content.
// Generates tmplsinit.gen.go.
func main() {
	err := genTemplatesVarFile()
	if err != nil {
		panic(err)
	}
}

func genTemplatesVarFile() error {
	m, err := makeTemplatesVar()
	if err != nil {
		return err
	}
	var bs []byte
	bs, err = makeTemplatesVarFile(m)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(templatesVarFileName, bs, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func makeTemplatesVar() (map[string]string, error) {
	m := make(map[string]string)
	dirs, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}
	var (
		fs      []os.FileInfo
		content []byte
	)
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
			m[string(fs[j].Name())] = string(content)
		}
	}
	return m, nil
}

func makeTemplatesVarFile(data map[string]string) ([]byte, error) {
	t := template.New("base")
	t.Parse(templatesVarFileSample)
	buf := bytes.NewBuffer(make([]byte, 0))
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
