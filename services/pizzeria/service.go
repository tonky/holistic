package pizzeria

import (
	"hf/holistic/domain"
	"hf/holistic/services"
	"io/fs"
	"log"
	"os"
	"text/template"

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
			"cap":      builtin.Capitalize,
			"port":     (*int)(nil),
			"handlers": &ps.Endpoints,
		},
	}

	// Build the program.
	temp, err := scriggo.BuildTemplate(fsys, st, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = temp.Run(os.Stdout, map[string]any{"port": 3001}, nil)
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
	getOrder := services.Endpoint{
		Name:   "orders",
		Method: services.Read,
		In: []services.InputParam{
			{Where: "path", What: domain.Object{Name: "food.OrderID"}, Validation: "required,len=16"},
		},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK:          "food.Order",
			services.ResponseNotFound:    "OrderNotFound",
			services.ResponseServerError: "http.ServerError",
		},
	}

	createOrder := services.Endpoint{
		Name:   "orders",
		Method: services.Create,
		In:     []services.InputParam{{Where: "body", What: domain.Object{Name: "food.Order"}}},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK:          "food.Order",
			services.ResponseServerError: "http.ServerError",
		},
	}

	return services.Service{Name: "pizzeria", T: services.HTTP, Endpoints: []services.Endpoint{getOrder, createOrder}}
}

type OrderNotFound struct{}
