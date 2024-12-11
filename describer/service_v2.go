package describer

import (
	"tonky/holistic/typs"

	"github.com/open2b/scriggo/builtin"
)

type ServiceV2 struct {
	Name           string
	Rpc            RPC
	Dependencies   Deps
	Endpoints      []EndpointV2
	Secrets        map[string]string
	ConfigItems    []ConfigItemV2
	AppConfigItems []ConfigItemV2
	Postgres       EndpointGroups
	KafkaProducers []TopicDesc2
	KafkaConsumers []TopicDesc2
	Clients        []InfraV2
	Logger         InfraInterface
	Tele           InfraInterface
}

type EndpointGroups []EndpointGroup

type EndpointV2 struct {
	Name   string
	In     typs.Object3
	Out    typs.Object3
	Errors []ErrorV2
}

type EndpointGroup struct {
	Name      string
	Endpoints []EndpointV2
}

func (eg EndpointGroup) InterfaceName() string {
	return "I" + eg.StructName()
}

func (eg EndpointGroup) StructName() string {
	return builtin.Capitalize(eg.Name)
}

type ErrorV2 struct {
	StatusCode int
	Message    string
}

type InfraV2 struct {
	Name  string
	Model typs.Object3
}

type InfraInterface struct {
	Interface typs.Object3
	Model     typs.Object3
}

type ConfigItemV2 struct {
	Model      typs.Object3
	SplitWords bool
	Default    string
}

func (egs EndpointGroups) AbsImports(ctx typs.Object3) []string {
	var imports []string

	for _, eg := range egs {
		for _, e := range eg.Endpoints {
			imports = append(imports, e.In.AbsImports(ctx)...)
			imports = append(imports, e.Out.AbsImports(ctx)...)
		}
	}

	return typs.ImportsDedup(imports)
}

func (s ServiceV2) AbsImports(ctx typs.Object3) []string {
	var imports []string

	imports = append(imports, s.Tele.Interface.AbsImports(ctx)...)
	// imports = append(imports, s.Logger.Model.AbsImports(ctx)...)

	for _, e := range s.Endpoints {
		imports = append(imports, e.In.AbsImports(ctx)...)
		imports = append(imports, e.Out.AbsImports(ctx)...)
	}

	return typs.ImportsDedup(imports)
}
