package describer

import (
	"strings"
	"tonky/holistic/typs"

	"github.com/open2b/scriggo/builtin"
)

type TopicDesc struct {
	Name         string
	TopicName    string
	DomainObject typs.Object
}

func (td TopicDesc) InterfaceName() string {
	return "I" + builtin.Capitalize(td.Name)
}

func (td TopicDesc) StructName() string {
	return builtin.Capitalize(td.Name)
}

func (td TopicDesc) ModelName() string {
	return td.DomainObject.Domain + "." + td.DomainObject.Name
}

func (td TopicDesc) SnakeFileName() string {
	return strings.Replace(td.TopicName, ".", "_", -1)
}

func (td TopicDesc) AppVarName() string {
	return td.StructName()
}