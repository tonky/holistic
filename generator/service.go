package generator

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"tonky/holistic/generator/services"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func GenService(s services.Service) {
	fmt.Printf("Generating %s service Go files\n", s.Name)

	template_dir := "templates"
	stn := "server_net_rpc.tpl"
	ctn := "client.tpl"
	service_config_tpl := "service_config.tpl"

	if err := os.MkdirAll(filepath.Join(".", "clients"), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "services", s.Name), os.ModePerm); err != nil {
		panic(err)
	}

	tplGenPath := map[string]string{
		stn:                fmt.Sprintf("services/%s/server_%s.go", s.Name, s.Rpc.String()),
		service_config_tpl: fmt.Sprintf("services/%s/config.go", s.Name),
		ctn:                fmt.Sprintf("clients/%s_client.go", s.Name),
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
			"handlers":     &s.Endpoints,
			"service_name": &s.Name,
			"config_items": &s.ConfigItems,
			"infra":        &s.Infra,
		},
	}

	writeTemplate(fsys, stn, opts, nil, tplGenPath[stn])
	writeTemplate(fsys, ctn, opts, nil, tplGenPath[ctn])
	writeTemplate(fsys, service_config_tpl, opts, nil, tplGenPath[service_config_tpl])

	fmt.Println("Generated client")
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