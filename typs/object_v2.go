package typs

import (
	"fmt"
	"strings"
)

type FieldType2 int

const (
	Int FieldType2 = iota
	Float
	String
	UUID
	DateTime
	Struct
	ObjectList
)

func (ft FieldType2) GoType() string {
	switch ft {
	case Int:
		return "int"
	case Float:
		return "float64"
	case String:
		return "string"
	case UUID:
		return "uuid.UUID"
	case DateTime:
		return "time.Time"
	case Struct:
		return "struct"
	}

	return "undefined"
}

type Object2 struct {
	Typ         FieldType2
	Domain      string
	Name        string
	Al          string
	Path        []string
	PackagePath string
	Fields      []Object2
}

func (o Object2) IsBuiltin() bool {
	return o.Typ == Int || o.Typ == Float || o.Typ == String
}

func (o Object2) IsPrimitive() bool {
	return o.Typ == UUID || o.Typ == DateTime
}

func (o Object2) TypeStr() string {
	switch o.Typ {
	case Struct:
		return "struct"
	case ObjectList:
		return "object_list"
	}

	return "basic"
}

func (o Object2) FQImport() string {
	if o.IsBuiltin() {
		return ""
	}

	if o.IsPrimitive() {
		switch o.Typ {
		case UUID:
			return "github.com/google/uuid"
		case DateTime:
			return "time"
		}
	}

	return fmt.Sprintf("%s/%s", o.PackagePath, o.FsRelPath())
}

func (o Object2) Alias() string {
	return o.Al
}

func (o Object2) GoName() string {
	if o.Al != "" {
		return o.Al
	}

	return o.Name
}

func (o Object2) FQType(ctx Object2) string {
	fmt.Printf(">> %s.%s FQType(%s.%s)\n", o.Domain, o.Name, ctx.Domain, ctx.Name)

	if o.Domain == ctx.Domain && o.Name != ctx.Name {
		return o.Name
	}

	return o.Typ.GoType()
}

type GoFQImport interface {
	FQPath() string
	Alias() string
}

type FQImport struct {
	Path string
	Al   string
}

func (fqi FQImport) FQPath() string {
	return fqi.Path
}

func (fqi FQImport) Alias() string {
	return fqi.Al
}

func (o Object2) GoImports(ctx Object2) []FQImport {
	fmt.Println(o.Name+".GoImports ", "for context: ", ctx.Name)
	res := []FQImport{}

	if o.IsPrimitive() {
		fmt.Println("Field is proimitive: ", o.Name, "adding dependency: ", o.Typ.GoType())

		return []FQImport{{Path: o.FQImport(), Al: ""}}
	}

	for _, t := range o.Fields {
		fmt.Println("Field: ", t.Name, "FQImport: ", t.FQImport())

		if t.Domain == ctx.Domain {
			fmt.Println("Same domain: ", t.Domain, ", skipping...")
			continue
		}

		if t.IsBuiltin() {
			fmt.Println("FQImports() is builtin: ", t.FQImport(), "skipping", t.Name)
			continue
		}

		fmt.Println("adding dependency: ", t.FQImport())
		res = append(res, FQImport{Path: t.FQImport(), Al: ""})
	}

	fmt.Println("Oject2.GoImports: ", o.Name, res)

	return res
}

func (o Object2) ShouldGenerate() bool {
	return o.Domain != "" && o.Domain != "hellofresh"
}

func (o Object2) FsRelPath() string {
	p := append([]string{"domain"}, o.Path...)
	p = append(p, o.Domain)

	return strings.Join(p, "/")
}
