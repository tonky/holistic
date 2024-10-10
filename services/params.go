package services

import (
	"fmt"
	"hf/holistic/domain"

	"github.com/open2b/scriggo/builtin"
)

type Transport string

const (
	HTTP Transport = "http"
	GRPC Transport = "grpc"
)

// topics: generate static list of topics from schema registry
type Topic string

// service names: generate static list of services
type ServiceName string

const (
	Orders ServiceName = "orders"
)

func (t Transport) String() string {
	if t == HTTP {
		return "HTTP"
	}

	return "gRPC"
}

type InputParam struct {
	Where      string
	What       domain.Object
	Validation string
}

func (ip InputParam) URLParamName() string {
	return ip.What.FieldName()
}

type Inputs []InputParam

func (i Inputs) Path() string {
	path := "/"

	for _, ip := range i {
		if ip.Where == "path" {
			path += fmt.Sprintf("{%s}", ip.What.Name)
		}
	}

	return path
}

func (i Inputs) FuncArgs() string {
	return "fa funcArgs"
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
	return string(e.Method) + builtin.Capitalize(e.Name)
}

func (e Endpoint) FuncArgs() string {
	return e.In.FuncArgs()
}

type Service struct {
	Name      string
	T         Transport
	Endpoints []Endpoint
	Secrets   map[string]string
	Publishes []Topic
	Consumes  []Topic
	// ACLs
}

func (e Endpoint) Debug(t Transport) string {
	method := ""
	out := ""

	if t == HTTP {
		if e.Method == Read {
			method = "GET"
		}
	}

	for rt, o := range e.Out {
		out += fmt.Sprintf("    %-9s: %s\n", rt, o)
	}

	return fmt.Sprintf("%s %s\n%s", method, e.In.Path(), out)
}

func (s Service) Debug() string {
	header := fmt.Sprintf("%s\n====\n", s.T)
	res := ""

	for _, e := range s.Endpoints {
		res += e.Debug(s.T)
	}

	return header + res
}
