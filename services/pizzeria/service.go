package pizzeria

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	stn := "server_net_rpc.tpl"
	ctn := "client.tpl"
	service_config_tpl := "service_config.tpl"

	ps := New()

	if err := os.MkdirAll(filepath.Join(".", "gen", "clients"), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "gen", "services", ps.Name), os.ModePerm); err != nil {
		panic(err)
	}

	tplGenPath := map[string]string{
		stn:                fmt.Sprintf("gen/services/%s/server_%s.go", ps.Name, ps.Rpc.String()),
		service_config_tpl: fmt.Sprintf("gen/services/%s/config.go", ps.Name),
		ctn:                fmt.Sprintf("gen/clients/%s_client.go", ps.Name),
	}

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

	writeTemplate(fsys, stn, opts, nil, tplGenPath[stn])
	writeTemplate(fsys, ctn, opts, nil, tplGenPath[ctn])
	writeTemplate(fsys, service_config_tpl, opts, nil, tplGenPath[service_config_tpl])

	fmt.Println("Generated client")
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
