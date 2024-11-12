package generator

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"tonky/holistic/typs"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func GenModels(models []typs.Object) {
	for _, model := range models {
		newpath := filepath.Join(".", "domain", model.Domain)

		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	st := "domain_model.tpl"

	fSys := os.DirFS("templates")

	contents, err := fs.ReadFile(fSys, st)
	if err != nil {
		log.Fatal(err)
	}

	fsys := scriggo.Files{st: contents}

	for _, model := range models {
		imports := model.GoImports()

		opts := &scriggo.BuildOptions{
			Globals: native.Declarations{
				"imports": &imports,
				"fields":  &model.Fields,
				"domain":  &model.Domain,
				"model":   &model.Name,
			},
		}

		// Build the program.
		temp, err := scriggo.BuildTemplate(fsys, st, opts)
		if err != nil {
			log.Fatal(err)
		}

		// open a file and get a writer
		fm, err := os.Create(fmt.Sprintf("./domain/%s/%s.go", model.Domain, builtin.ToLower(model.Name)))
		if err != nil {
			log.Fatal(err)
		}

		err = temp.Run(fm, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
