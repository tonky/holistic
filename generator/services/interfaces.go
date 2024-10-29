package services

import (
	"fmt"
	"strings"
)

type JustInterface struct {
	Name    string
	Struct  string
	Deps    InterfaceDeps
	Methods []InterfaceMethod
}

func (ji JustInterface) InterfaceName() string {
	if ji.Name != "" {
		return ji.Name
	}

	return "I" + ji.Struct
}

func (ji JustInterface) Imports(s Service) []ClientImport {
	var out []ClientImport

	for _, m := range ji.Methods {
		out = append(out, m.Arg.Import(s))
		out = append(out, m.Ret.Import(s))
	}

	return out
}

type InterfaceDeps map[string]FQImport

func (id InterfaceDeps) StructArgs() []InfraObject {
	out := []InfraObject{}

	for k, fqi := range id {
		out = append(out, InfraObject{Name: k, Typ: fqi.FQModel()})
	}

	return out
}

func (id InterfaceDeps) StructArgsStr() string {
	var args []string

	for varName, fqi := range id {
		args = append(args, fmt.Sprintf("%s %s", varName, fqi.FQModel()))
	}

	return strings.Join(args, ", ")
}

type FQImport struct {
	Package string
	Model   string
	RelPath string
	AbsPath string
}

func (fqi FQImport) FQModel() string {
	split := strings.Split(fqi.RelPath, "/")

	return split[len(split)-1] + "." + fqi.Model
}

func (fqi FQImport) FQImport(s Service) string {
	if fqi.AbsPath != "" {
		return fqi.AbsPath
	}

	if fqi.Package == "app" {
		return fmt.Sprintf("apps/%s/%s", strings.ToLower(s.Name), fqi.RelPath)
	}

	if fqi.Package == "svc" {
		return fmt.Sprintf("services/%s/%s", strings.ToLower(s.Name), fqi.RelPath)
	}

	return fqi.RelPath

}

func NewFQImport(pkg, model, pkgPath, repo string) FQImport {
	return FQImport{
		Package: pkg,
		Model:   model,
		RelPath: pkgPath,
	}
}

func ParseDep(depPath string) FQImport {
	if depPath == "" {
		panic("empty dep path in ParseDep")
	}

	// absolute path
	if strings.Contains(depPath, "/") {
		split := strings.Split(depPath, ".")

		pSplit := strings.Split(split[0], "/")

		return FQImport{
			Package: pSplit[len(split)-1],
			Model:   split[len(split)-1],
			AbsPath: split[0],
		}
	}

	split := strings.Split(depPath, ".")

	if len(split) <= 1 {
		panic("dependency import must have at least 2 parts with '.', got: " + depPath)
	}

	// app model dependency?
	if len(split) == 2 {
		return FQImport{
			Package: split[0],
			Model:   split[1],
		}
	}

	return FQImport{
		Package: split[len(split)-2],
		Model:   split[len(split)-1],
		RelPath: strings.Join(split[0:len(split)-2], "/"),
	}
}