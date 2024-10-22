package domain

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

type Kind int

const (
	Int Kind = iota
	Float
	String
	UUID
	StringList
	ObjectList
)

type FieldType string

func (ft FieldType) StructType(domain string) string {
	fms := strings.Split(string(ft), ".")

	if len(fms) != 2 {
		return string(ft)
	}

	domainName := fms[0]
	modelName := fms[1]

	if domainName == domain {
		return modelName
	}

	return string(ft)
}

type Object struct {
	Name   string
	Fields []Field
}

func (o Object) GoImports(domainName string) []string {
	res := []string{}

	for _, t := range o.Fields {
		tdn := strings.Split(string(t.T), ".")

		if len(tdn) != 2 {
			continue
		}

		typeDomain := tdn[0]

		if domainName == typeDomain {
			continue
		}

		domainImport, err := resolveDomainImport(typeDomain)
		if err != nil {
			panic(err)
		}

		res = append(res, domainImport)

	}

	return res
}

func resolveDomainImport(d string) (string, error) {
	switch d {
	case "billing":
		return "tonky/holistic/domain/billing", nil
	case "uuid":
		return "github.com/google/uuid", nil
	}

	return "undefined", fmt.Errorf("domain import unresolved")
}

func Generate() {
	domainName := "food"

	newpath := filepath.Join(".", "gen", "domain", domainName)
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	st := "domain_model.tpl"

	fSys := os.DirFS("templates")

	contents, err := fs.ReadFile(fSys, st)
	if err != nil {
		log.Fatal(err)
	}

	fsys := scriggo.Files{st: contents}

	for _, foodModel := range foods {
		imports := foodModel.GoImports(domainName)

		opts := &scriggo.BuildOptions{
			Globals: native.Declarations{
				"imports": &imports,
				"fields":  &foodModel.Fields,
				"domain":  &domainName,
				"model":   &foodModel.Name,
			},
		}

		// Build the program.
		temp, err := scriggo.BuildTemplate(fsys, st, opts)
		if err != nil {
			log.Fatal(err)
		}

		// open a file and get a writer
		fm, err := os.Create(fmt.Sprintf("./domain/%s/%s.go", domainName, builtin.ToLower(foodModel.Name)))
		if err != nil {
			log.Fatal(err)
		}

		err = temp.Run(fm, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (o Object) FieldName() string {
	tmp := builtin.Split(o.Name, ".")

	return tmp[len(tmp)-1]
}

type Field struct {
	Name string
	T    FieldType
}

var FoodOrderID = Object{
	Name: "OrderID",
	Fields: []Field{
		{Name: "ID", T: "uuid.UUID"},
	},
}

var FoodOrder = Object{
	Name: "Order",
	Fields: []Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Content", T: "string"},
	},
}

var foods = []Object{FoodOrderID, FoodOrder}
