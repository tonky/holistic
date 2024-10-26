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
	app_tpl := "app.tpl"
	repo_tpl := "repository.tpl"
	kafka_prod_tpl := "kafka_producer.tpl"

	if err := os.MkdirAll(filepath.Join(".", "clients"), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "services", s.Name), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "apps", s.Name), os.ModePerm); err != nil {
		panic(err)
	}

	tplGenPath := map[string]string{
		stn:                fmt.Sprintf("services/%s/server_%s.go", s.Name, s.Rpc.String()),
		service_config_tpl: fmt.Sprintf("services/%s/config.go", s.Name),
		ctn:                fmt.Sprintf("clients/%s_client.go", s.Name),
		app_tpl:            fmt.Sprintf("apps/%s/gen_%s_app.go", s.Name, s.Name),
	}

	fsys := scriggo.Files{
		stn:                readContent(template_dir, stn),
		ctn:                readContent(template_dir, ctn),
		service_config_tpl: readContent(template_dir, service_config_tpl),
		app_tpl:            readContent(template_dir, app_tpl),
		repo_tpl:           readContent(template_dir, repo_tpl),
		kafka_prod_tpl:     readContent(template_dir, kafka_prod_tpl),
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

	appDeps := []AppDep{}
	opts.Globals["app_deps"] = &appDeps

	for _, pg := range s.Postgres {
		appDeps = append(appDeps, pg)
		opts.Globals["repo"] = &pg
		opts.Globals["kind"] = "postgres"

		outFile := fmt.Sprintf("apps/%s/gen_%s_repository_postgres.go", s.Name, pg.Name)
		writeTemplate(fsys, repo_tpl, opts, nil, outFile)
	}

	for _, kp := range s.KafkaProducers {
		appDeps = append(appDeps, kp)
		opts.Globals["kp"] = &kp

		outFile := fmt.Sprintf("apps/%s/gen_%s_producer_kafka.go", s.Name, kp.Name)
		writeTemplate(fsys, kafka_prod_tpl, opts, nil, outFile)
	}

	writeTemplate(fsys, stn, opts, nil, tplGenPath[stn])
	writeTemplate(fsys, ctn, opts, nil, tplGenPath[ctn])
	writeTemplate(fsys, service_config_tpl, opts, nil, tplGenPath[service_config_tpl])
	writeTemplate(fsys, app_tpl, opts, nil, tplGenPath[app_tpl])

	fmt.Println("Generated Go files for service", s.Name)
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

type AppDep interface {
	AppVarName() string
	InterfaceName() string
}
