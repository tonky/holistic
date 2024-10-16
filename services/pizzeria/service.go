package pizzeria

import (
	"fmt"
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
	fmt.Println("Generated pizza service Go files")

	stn := "server_http.tpl"
	ctn := "client.tpl"

	fSys := os.DirFS("templates")

	sContents, err := fs.ReadFile(fSys, stn)
	if err != nil {
		log.Fatal(err)
	}

	cContents, err := fs.ReadFile(fSys, ctn)
	if err != nil {
		log.Fatal(err)
	}

	ps := New()

	fsys := scriggo.Files{
		stn: sContents,
		ctn: cContents,
	}

	opts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"cap":          builtin.Capitalize,
			"port":         (*int)(nil),
			"handlers":     &ps.Endpoints,
			"service_name": ps.Name,
		},
	}

	sTemp, err := scriggo.BuildTemplate(fsys, stn, opts)
	if err != nil {
		log.Fatal(err)
	}

	sGoFile, err := os.Create("./gen/services/pizzeria/http/server_http.go")
	if err != nil {
		log.Fatal(err)
	}

	err = sTemp.Run(sGoFile, map[string]any{"port": 3001}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated server")

	clientOpts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"cap":          builtin.Capitalize,
			"handlers":     &ps.Endpoints,
			"service_name": ps.Name,
		},
	}

	cGoFile, err := os.Create("./gen/clients/pizzeria_client.go")
	if err != nil {
		log.Fatal(err)
	}

	cTemp, err := scriggo.BuildTemplate(fsys, ctn, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	err = cTemp.Run(cGoFile, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated client")
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
