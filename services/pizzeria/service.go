package pizzeria

import (
	"io/fs"
	"log"
	"os"
	"text/template"
	"tonky/holistic/services"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

type ServiceTpl struct {
	Port int
}

func GenScrig() {
	st := "server_http.tpl"

	fSys := os.DirFS("templates")

	contents, err := fs.ReadFile(fSys, st)
	if err != nil {
		log.Fatal(err)
	}

	ps := New()

	fsys := scriggo.Files{"server_http.tpl": contents}

	opts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"cap":          builtin.Capitalize,
			"port":         (*int)(nil),
			"handlers":     &ps.Endpoints,
			"service_name": ps.Name,
		},
	}

	// Build the program.
	temp, err := scriggo.BuildTemplate(fsys, st, opts)
	if err != nil {
		log.Fatal(err)
	}

	// open a file and get a writer
	fm, err := os.Create("./gen/services/pizzeria/http/server_http.go")
	if err != nil {
		log.Fatal(err)
	}

	err = temp.Run(fm, map[string]any{"port": 3001}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Generate() {
	svc := ServiceTpl{3000}
	tmplFile := "templates/server_http.tpl"

	tmpl, err := template.New("server_http.tpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, svc)
	if err != nil {
		panic(err)
	}
}

func New() services.Service {
	// rpc: net/rpc, twirp
	getOrder := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK:          "food.Order",
			services.ResponseNotFound:    "OrderNotFound",
			services.ResponseServerError: "http.ServerError",
		},
	}

	createOrder := services.Endpoint{
		Name:   "order",
		Method: services.Create,
		In:     services.Inputs{Name: "food.Order"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK:          "food.Order",
			services.ResponseServerError: "http.ServerError",
		},
	}

	return services.Service{Name: "pizzeria", Rpc: services.GoNative, Endpoints: []services.Endpoint{getOrder, createOrder}}
}

type OrderNotFound struct{}
