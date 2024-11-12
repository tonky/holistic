package typs

import (
	"fmt"
	"strings"

	"github.com/open2b/scriggo/builtin"
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

type Field struct {
	Name string
	T    FieldType
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

func (o Object) FieldName() string {
	tmp := builtin.Split(o.Name, ".")

	return tmp[len(tmp)-1]
}
