package services

import (
	"fmt"
	"strings"
	"tonky/holistic/generator/domain"

	"github.com/open2b/scriggo/builtin"
)

type RPC string

const (
	GoNative RPC = "net_rpc"
	Twirp    RPC = "twirp"
	GRPC     RPC = "grpc"
)

type ObjectType int

const (
	DomainType ObjectType = iota
	ServiceType
)

// topics: generate static list of topics from schema registry
type Topic string

// service names: generate static list of services
type ServiceName string

const (
	Orders ServiceName = "orders"
)

func (r RPC) String() string {
	if r == GoNative {
		return "net_rpc"
	}

	return "twirp"
}

type InputParam struct {
	Where      string
	What       domain.Object
	Validation string
}

func (ip InputParam) URLParamName() string {
	return ip.What.FieldName()
}

type Inputs struct {
	// Typ  ObjectType
	Name string
	// Validation string
}

func (i Inputs) ModelName() string {
	return "modelName"
}

func (i Inputs) ServiceModel() string {
	if strings.Contains(i.Name, ".") {
		return i.Name
	}

	return "svc." + i.Name
}

func (i Inputs) SvcToApp() string {
	if strings.Contains(i.Name, ".") {
		return ""
	}

	return ".ToApp()"
}

func (i Inputs) String() string {
	return i.Name
}

type ResponseObject string

type ResponseType string

const (
	ResponseOK          ResponseType = "ok"
	ResponseNotFound    ResponseType = "not_found"
	ResponseServerError ResponseType = "server_error"
)

func (rt ResponseType) String() string {
	switch rt {
	case ResponseNotFound:
		return "NotFound"
	case ResponseOK:
		return "OK"
	case ResponseServerError:
		return "ServerError"
	}

	return "undefined"
}

type Endpoint struct {
	Name   string
	Method Method
	In     Inputs
	Out    map[ResponseType]ResponseObject
}

func (e Endpoint) FuncName() string {
	return builtin.Capitalize(string(e.Method) + builtin.Capitalize(e.Name))
}

type Service struct {
	Name           string
	Rpc            RPC
	Endpoints      []Endpoint
	Secrets        map[string]string
	Publishes      []Topic
	Consumes       []Topic
	ConfigItems    []ConfigItem
	Infra          []Infra
	Postgres       []Postgres
	KafkaProducers []KafkaProducer
	// specific infra
	// generic escape hatches
	// ACLs
}

func (e Endpoint) Debug(r RPC) string {
	out := ""

	for rt, o := range e.Out {
		out += fmt.Sprintf("    %-9s: %s\n", rt, o)
	}

	return fmt.Sprintf("%s %s\n%s", e.Method, e.In.String(), out)
}

func (s Service) Debug() string {
	header := fmt.Sprintf("%s\n====\n", s.Rpc)
	res := ""

	for _, e := range s.Endpoints {
		res += e.Debug(s.Rpc)
	}

	return header + res
}

type ConfigItem struct {
	Name string
	Typ  string
}

type InfraObject struct {
	Name string
	Typ  string
}

type InOut struct {
	Name string
	In   InfraObject
	Out  InfraObject
}

type Infra struct {
	Name  string
	Typ   string
	InOut []InOut
}

func (i Infra) ConfigVar() string {
	switch i.Typ {
	case "postgres":
		return "PostgresConfig"
	case "kafka":
		return "KafkaConfig"
	default:
		panic("Infra.ConfigVar(): unknown infra " + i.Typ)
	}
}

func (i Infra) AppVarName() string {
	return builtin.ToLower(i.Name) + "Repo"
}

func (i Infra) InterfaceName() string {
	return i.Name + "Repository"
}

func (i Infra) ImplName() string {
	return builtin.Capitalize(i.Typ) + i.Name
}

func (i Infra) ClientName() string {
	return "infra.New" + builtin.Capitalize(i.Typ) + "Client()"
}

func (i Infra) ConfigFQN() string {
	return "infra." + builtin.Capitalize(i.Typ) + "Config"
}

func (i Infra) ClientFQN() string {
	return "infra.New" + builtin.Capitalize(i.Typ) + "Client"
}

func (i Infra) ClientType() string {
	if i.Typ != "kafka" {
		return "Client"
	}

	if i.InOut[0].In.Name != "" {
		return "Consumer"
	}

	return "Producer"
}

func (i Infra) TopicName() string {
	if i.InOut[0].In.Name != "" {
		return i.InOut[0].In.Name
	}

	return i.InOut[0].Out.Name
}
