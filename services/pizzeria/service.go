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

	template_dir := "templates"
	stn := "server_http.tpl"
	ctn := "client.tpl"
	service_config_tpl := "service_config.tpl"

	ps := New()

	fsys := scriggo.Files{
		stn:                readContent(template_dir, stn),
		ctn:                readContent(template_dir, ctn),
		service_config_tpl: readContent(template_dir, service_config_tpl),
	}

	opts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"cap":          builtin.Capitalize,
			"port":         (*int)(nil),
			"handlers":     &ps.Endpoints,
			"service_name": &ps.Name,
			"config_items": &ps.ConfigItems,
			"infra":        &ps.Infra,
		},
	}

	writeTemplate(fsys, stn, opts, nil, "gen/services/pizzeria/http/server_http.go")
	writeTemplate(fsys, ctn, opts, nil, "gen/clients/pizzeria_client.go")
	writeTemplate(fsys, service_config_tpl, opts, nil, "gen/services/pizzeria/pizzeria_config.go")

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
			services.ResponseOK: "food.Order",
		},
	}

	createOrder := services.Endpoint{
		Name:   "order",
		Method: services.Create,
		In:     services.Inputs{Name: "food.Order"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "food.Order",
		},
	}

	return services.Service{
		Name:        "pizzeria",
		Rpc:         services.GoNative,
		Endpoints:   []services.Endpoint{getOrder, createOrder},
		ConfigItems: []services.ConfigItem{{Name: "ShouldMockApp", Typ: "bool"}},
		Infra:       []services.Infra{{Name: "postgres"}},
	}
}

func writeTemplate(fsys scriggo.Files, tplName string, opts *scriggo.BuildOptions, vars map[string]any, outFile string) {
	goFile, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}

	cTemp, err := scriggo.BuildTemplate(fsys, tplName, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = cTemp.Run(goFile, vars, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func readContent(dir string, file string) []byte {
	content, err := fs.ReadFile(os.DirFS(dir), file)
	if err != nil {
		log.Fatal(err)
	}

	return content
}
