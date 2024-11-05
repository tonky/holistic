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
	Domain string
	Name   string
	Fields []Field
}

func (o Object) GoImport() string {
	dimp, err := resolveDomainImport(o.Domain)
	if err != nil {
		panic(err)
	}

	return dimp
}

func (o Object) GoImports() []string {
	res := []string{}

	for _, t := range o.Fields {
		tdn := strings.Split(string(t.T), ".")

		if len(tdn) != 2 {
			continue
		}

		typeDomain := tdn[0]

		if o.Domain == typeDomain {
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
	case "food":
		return "tonky/holistic/domain/food", nil
	case "accounting":
		return "tonky/holistic/domain/accounting", nil
	case "shipping":
		return "tonky/holistic/domain/shipping", nil
	case "pricing":
		return "tonky/holistic/domain/pricing", nil
	case "uuid":
		return "github.com/google/uuid", nil
	case "time":
		return "time", nil
	}

	return "undefined", fmt.Errorf("domain import unresolved: %s", d)
}

func Generate() {
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

func (o Object) FieldName() string {
	tmp := builtin.Split(o.Name, ".")

	return tmp[len(tmp)-1]
}

type Field struct {
	Name string
	T    FieldType
}

var FoodOrderID = Object{
	Domain: "food",
	Name:   "OrderID",
	Fields: []Field{
		{Name: "ID", T: "uuid.UUID"},
	},
}

var FoodOrder = Object{
	Domain: "food",
	Name:   "Order",
	Fields: []Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Content", T: "string"},
		{Name: "IsFinal", T: "bool"},
	},
}

var AccountingOrder = Object{
	Domain: "accounting",
	Name:   "Order",
	Fields: []Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Cost", T: "int"},
	},
}

var ShippedOrder = Object{
	Domain: "shipping",
	Name:   "Order",
	Fields: []Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "ShippedAt", T: "time.Time"},
	},
}

var OrderPrice = Object{
	Domain: "pricing",
	Name:   "OrderPrice",
	Fields: []Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Cost", T: "int"},
	},
}

var models = []Object{FoodOrderID, FoodOrder, AccountingOrder, ShippedOrder, OrderPrice}
