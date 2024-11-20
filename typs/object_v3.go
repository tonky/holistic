package typs

import (
	"fmt"
	"log/slog"
	"strings"
)

type ObjectType string

const (
	Bool2     ObjectType = "bool"
	Int2      ObjectType = "int"
	Float2    ObjectType = "float"
	String2   ObjectType = "string"
	Struct2   ObjectType = "struct"
	UUID2     ObjectType = "uuid"
	Time2     ObjectType = "time"
	Duration2 ObjectType = "duration"
)

func (o Object3) ShouldGenerate() bool {
	return o.Kind == KindDomain
}

func (o Object3) GoType() string {
	if o.Kind == KindExternal {
		return o.Package() + "." + o.Name
	}

	switch o.Typ {
	case Bool2:
		return "bool"
	case Int2:
		return "int"
	case Float2:
		return "float64"
	case String2:
		return "string"
	case Struct2:
		return "struct"
	case UUID2:
		return "uuid.UUID"
	case Time2:
		return "time.Time"
	case Duration2:
		return "time.Duration"
	}

	return "GoType(): undefined"
}

func (o Object3) IsDomain() bool {
	return o.Kind == KindDomain
}

func (o Object3) IsClient() bool {
	return o.Kind == KindClient
}

func (o Object3) IsBuiltin() bool {
	return o.Typ == Int2 || o.Typ == Float2 || o.Typ == String2 || o.Typ == Bool2
}

func (o Object3) IsBasic() bool {
	return o.Typ == UUID2 || o.Typ == Time2 || o.Typ == Duration2
}

func (o Object3) GoQualifiedModel() string {
	slog.Info("GoQualifiedModel", slog.String("name", fmt.Sprintf("%s.%s", o.RelPath(), o.Name)))

	return fmt.Sprintf("%s.%s", o.Package(), o.Name)
}

func (o Object3) AbsPath() string {
	return fmt.Sprintf("%s/%s", o.Module, o.RelPath())
}

func (o Object3) GoStructModel(ctx Object3) string {
	slog.Info("GoStructModel", slog.String("name", o.Name), slog.String("context", ctx.Name))

	if ctx.IsClient() {
		return "svc." + o.Name

	}

	if o.IsDomain() {
		slog.Info("IsDomain()")

		if o.Module == ctx.Module && o.RelPath() == ctx.RelPath() {
			slog.Info("same abspath(), returning name")

			return o.Name
		}

		slog.Info("not the same abspath(), returning qualified model")

		return o.GoQualifiedModel()
	}

	slog.Info("not domain, returning GoType()")

	return o.GoType()
}

func (o Object3) BasicGoImport() string {
	slog.Info("BasicGoImport", slog.Any("object", o))

	if o.Kind == KindExternal {
		return o.AbsPath()
	}

	switch o.Typ {
	case UUID2:
		return "github.com/google/uuid"
	case Time2:
		return "time"
	case Duration2:
		return "time"
	}

	return "BasicGoImport(): undefined"
}

type ObjectKind int

const (
	KindDomain ObjectKind = iota
	KindBuiltIn
	KindBasic
	KindExternal
	KindClient
)

type Object3 struct {
	Typ          ObjectType
	Kind         ObjectKind
	Name         string
	ImportAlias  string
	RelativePath []string
	Module       string
	Fields       []Object3
}

func (o Object3) Package() string {
	if len(o.RelativePath) != 0 {
		return o.RelativePath[len(o.RelativePath)-1]
	}

	modPath := strings.Split(o.Module, "/")

	return modPath[len(modPath)-1]
}

func (o Object3) GoFieldName() string {
	if o.ImportAlias != "" {
		return o.ImportAlias
	}

	return o.Name
}

func (o Object3) AbsImports(ctx Object3) []string {
	return ImportsDedup(o.AbsImportsAll(ctx))
}

func (o Object3) AbsImportsAll(ctx Object3) []string {
	slog.Info("AbsImport", slog.String("name", fmt.Sprintf("%s.%s", o.RelPath(), o.Name)), slog.String("context", ctx.Name))
	// fmt.Println("> AbsImport", o.Name)

	if o.IsBuiltin() {
		slog.Info("..builtin or service, no import required")

		return []string{}
	}

	slog.Info("..not builtin. Package: ", slog.String("package", o.Module), slog.String("relpath", o.RelPath()))
	slog.Info("..not builtin. Context: ", slog.String("package", ctx.Module), slog.String("relpath", ctx.RelPath()))

	if o.Module == "" {
		slog.Info("....no module")
		return []string{o.BasicGoImport()}
	}

	if o.Module == ctx.Module && o.RelPath() == ctx.RelPath() {
		slog.Info("....same package and relpath")

		if o.Name != ctx.Name { // same import, different object
			slog.Info("......different object")
			return []string{}
		}

		if o.IsBasic() {
			slog.Info("......basic", slog.Any("name", o.Name))
			return []string{o.BasicGoImport()}
		}

		slog.Info("....getting field imports")

		// self-imports for complex object - all fields imports
		fieldImports := []string{}

		for _, field := range o.Fields {
			imps := field.AbsImportsAll(o)

			fieldImports = append(fieldImports, imps...)
		}

		return fieldImports
	}

	slog.Info("returning full import")

	return []string{fmt.Sprintf("%s/%s", o.Module, o.RelPath())}
}

func (o Object3) RelPath() string {
	return strings.Join(o.RelativePath, "/")
}

func ImportsDedup(imports []string) []string {
	seen := map[string]bool{}

	res := []string{}

	for _, imp := range imports {
		if ok := seen[imp]; ok {
			continue
		}

		seen[imp] = true

		res = append(res, imp)
	}

	return res
}
