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
	HTTP     RPC = "http"
)

type Deps string

const (
	Struct   Deps = "plain_struct"
	SamberDO Deps = "samber_do"
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

	if r == HTTP {
		return "http"
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
	Method MethodAction
	In     Inputs
	Out    map[ResponseType]ResponseObject
}

func (e Endpoint) FuncName() string {
	return builtin.Capitalize(string(e.Method) + builtin.Capitalize(e.Name))
}

type Service struct {
	Name           string
	Rpc            RPC
	Dependencies   Deps
	Endpoints      []Endpoint
	Secrets        map[string]string
	Publishes      []Topic
	Consumes       []Topic
	ConfigItems    []ConfigItem
	AppConfigItems []ConfigItem
	Infra          []Infra
	Interfaces     []JustInterface
	Postgres       []Postgres
	KafkaProducers []TopicDesc
	KafkaConsumers []TopicDesc
	Clients        []Client
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

func (s Service) ClientImports() []ClientImport {
	res := []ClientImport{}

	for _, e := range s.Endpoints {
		split := strings.Split(e.In.Name, ".")

		fmt.Println("split: ", split)

		if len(split) == 2 {
			res = append(res, ClientImport{RelPath: "domain/" + split[0]})
		} else {
			res = append(res, ClientImport{RelPath: "services/" + s.Name, Alias: "svc"})
		}

		splitOut := strings.Split(string(e.Out[ResponseOK]), ".")

		fmt.Println("split out: ", splitOut)

		if len(splitOut) == 2 {
			res = append(res, ClientImport{RelPath: "domain/" + splitOut[0]})
		} else {
			res = append(res, ClientImport{RelPath: "services/" + s.Name, Alias: "svc"})
		}
	}

	fmt.Printf("ClientImports: %+v\n", res)

	seen := map[string]bool{}
	dedup := []ClientImport{}

	for _, ci := range res {
		if seen[ci.RelPath] {
			continue
		}

		seen[ci.RelPath] = true

		dedup = append(dedup, ci)
	}

	fmt.Printf("ClientImports dedup: %+v\n", dedup)

	return dedup
}

type ClientImport struct {
	Alias   string
	RelPath string
}

func (ci ClientImport) String(mod string) string {
	if ci.Alias != "" {
		return fmt.Sprintf("%s \"%s/%s\"", ci.Alias, mod, ci.RelPath)
	}

	return fmt.Sprintf("\"%s/%s\"", mod, ci.RelPath)
}

type ConfigItem struct {
	Name       string
	Typ        string
	SplitWords bool
	Default    string
}

type InfraObject struct {
	Name string
	Typ  string
}

func (io InfraObject) Import(s Service) (*ClientImport, error) {
	split := strings.Split(io.Typ, ".")

	// app model, no need to import
	if len(split) <= 1 {
		err := fmt.Errorf("InfraObject.Import: %+v not enough data in %s", s, io.Typ)

		fmt.Println(err)

		return nil, err
	}

	// not enough data, assume 'domain'?
	if len(split) == 2 {
		return &ClientImport{RelPath: strings.Join([]string{"domain", split[0]}, "/")}, nil
	}

	return &ClientImport{RelPath: strings.Join(split[0:len(split)-3], "/")}, nil
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
	return "I" + i.Name
}

func (i Infra) StructName() string {
	return i.AppVarName()
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

type Client struct {
	VarName string
	IName   string
}

func (c Client) AppVarName() string {
	return c.VarName
}

func (c Client) InterfaceName() string {
	return c.IName
}

func (c Client) StructName() string {
	return c.VarName
}

func (c Client) PackageName() string {
	return "clients"
}

func (c Client) AppImportPackageName() string {
	return c.PackageName()
}

func (c Client) ConfigVarName() string {
	return c.StructName()
}

func (c Client) ConfigVarType() string {
	return "clients.Config"
}
