package main

import (
	"github.com/ymz-ncnk/dvargen"
)

func main() {
	var (
		conf = func() dvargen.Conf {
			conf := dvargen.DefConf
			conf.Path = "text_template"
			return conf
		}()
		vDesc = dvargen.DVarDesc{
			Dir:     "text_template/templates/go/",
			Varname: "templates",
			Package: "text_template",
		}
	)
	err := dvargen.New().GenerateAs(vDesc, conf)
	if err != nil {
		panic(err)
	}
}
