package generator

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"tonky/holistic/typs"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func (g ServiceGen) GenModel2(model typs.Object2) {
	fmt.Printf("Generating domain object Go files: %s\n", model.Name)

	if !model.ShouldGenerate() {
		fmt.Println("...skipping")
		return
	}

	// recursively generate models for fields
	for _, field := range model.Fields {
		g.GenModel2(field)
	}

	domainDirPath := model.FsRelPath()

	fmt.Println("Creating model directory: ", domainDirPath)

	err := os.MkdirAll(domainDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	st := "domain_model2.tpl"

	contents, err := fs.ReadFile(os.DirFS(g.TemplatePath), st)
	if err != nil {
		log.Fatal(err)
	}

	fsys := scriggo.Files{st: contents}

	opts := &scriggo.BuildOptions{Globals: native.Declarations{"model": &model}}

	// Build the program.
	temp, err := scriggo.BuildTemplate(fsys, st, opts)
	if err != nil {
		log.Fatal("can't build template: ", err)
	}

	// open a file and get a writer
	fm, err := os.Create(fmt.Sprintf("./%s/%s.go", model.FsRelPath(), model.Name))
	if err != nil {
		log.Fatal("can't create file: ", err)
	}

	err = temp.Run(fm, nil, nil)
	if err != nil {
		log.Fatal("can't run template: ", err)
	}
}

func (g ServiceGen) GenModel3(model typs.Object3) error {
	slog.Info("Generating domain object v3 Go files", slog.Any("model", model))

	if !model.ShouldGenerate() {
		fmt.Println("...skipping")
		return nil
	}

	// recursively generate models for fields
	for _, field := range model.Fields {
		g.GenModel3(field)
	}

	domainDirPath := model.RelPath()

	fmt.Println("Creating model directory: ", domainDirPath)

	err := os.MkdirAll(domainDirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create directory: %w", err)
	}

	// open a file and get a writer
	fm, err := os.Create(fmt.Sprintf("./%s/%s.go", model.RelPath(), model.Name))
	if err != nil {
		return fmt.Errorf("can't create file: %w", err)
	}

	return g.WriteModel2(model, fm)
}

func (g ServiceGen) WriteModel2(model typs.Object3, w io.Writer) error {
	slog.Info("Writing domain object v3 Go vcode: %s\n", slog.String("model", model.Name))

	st := "domain_model2.tpl"

	contents, err := fs.ReadFile(os.DirFS(g.TemplatePath), st)
	if err != nil {
		return err
	}

	fsys := scriggo.Files{st: contents}

	opts := &scriggo.BuildOptions{Globals: native.Declarations{"model": &model}}

	temp, err := scriggo.BuildTemplate(fsys, st, opts)
	if err != nil {
		return err
	}

	return temp.Run(w, nil, nil)
}

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
		fm, err := os.Create(fmt.Sprintf("./domain/%s/gen_%s.go", model.Domain, builtin.ToLower(model.Name)))
		if err != nil {
			log.Fatal(err)
		}

		err = temp.Run(fm, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
