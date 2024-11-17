package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"tonky/holistic/describer"
	"tonky/holistic/typs"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func (g ServiceGen) Generate2(s describer.ServiceV2) error {
	fmt.Printf("Generating %s service Go files\n", s.Name)

	template_dir := g.TemplatePath

	service_net_rpc_tpl := "server_net_rpc_v2.tpl"
	service_http_tpl := "server_http_v2.tpl"
	client_net_rpc := "client_net_rpc_v2.tpl"
	client_http_tpl := "client_http_v2.tpl"
	service_config_tpl := "service_config_v2.tpl"
	app_config_tpl := "app_config_v2.tpl"
	app_tpl := "app_v2.tpl"
	repo_pg_tpl := "repository_postgres_v2.tpl"

	if err := os.MkdirAll(filepath.Join(".", "clients", s.Name+"Client"), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "services", s.Name), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(".", "apps", s.Name), os.ModePerm); err != nil {
		panic(err)
	}

	tplGenPath := map[string]string{
		service_net_rpc_tpl: fmt.Sprintf("services/%s/server_%s.go", s.Name, s.Rpc.String()),
		service_http_tpl:    fmt.Sprintf("services/%s/server_%s.go", s.Name, s.Rpc.String()),
		service_config_tpl:  fmt.Sprintf("services/%s/config.go", s.Name),
		client_net_rpc:      fmt.Sprintf("clients/%sClient/gen_client_%s.go", s.Name, s.Rpc.String()),
		client_http_tpl:     fmt.Sprintf("clients/%sClient/gen_client_%s.go", s.Name, s.Rpc.String()),
		app_tpl:             fmt.Sprintf("apps/%s/gen_%s_app.go", s.Name, s.Name),
		app_config_tpl:      fmt.Sprintf("apps/%s/gen_config.go", s.Name),
	}

	fsys := scriggo.Files{
		service_net_rpc_tpl: readContent(template_dir, service_net_rpc_tpl),
		service_http_tpl:    readContent(template_dir, service_http_tpl),
		client_net_rpc:      readContent(template_dir, client_net_rpc),
		client_http_tpl:     readContent(template_dir, client_http_tpl),
		service_config_tpl:  readContent(template_dir, service_config_tpl),
		app_tpl:             readContent(template_dir, app_tpl),
		app_config_tpl:      readContent(template_dir, app_config_tpl),
		repo_pg_tpl:         readContent(template_dir, repo_pg_tpl),
	}

	opts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"modulePath": "tonky/holistic",
			"service":    &s,
			"cap":        builtin.Capitalize,
			"ctx":        &typs.Object3{},
		},
	}

	appDeps := AppDeps{}
	infraDeps := AppDeps{}
	interfaces := AppDeps{}
	clientDeps := AppDeps{}
	// configImports := s.ClientImports()

	opts.Globals["app_deps"] = &appDeps
	opts.Globals["infra_deps"] = &infraDeps
	opts.Globals["interfaces"] = &interfaces
	opts.Globals["client_deps"] = &clientDeps
	// opts.Globals["client_relative_imports"] = &configImports

	for _, pg := range s.Postgres {
		opts.Globals["repo"] = &pg

		outFile := fmt.Sprintf("apps/%s/gen_%s_repository_postgres_v2.go", s.Name, pg.Name)
		writeTemplate(fsys, repo_pg_tpl, opts, nil, outFile)

		// generate repo models
		for _, eg := range pg.Endpoints {
			g.GenModel3(eg.In)
			g.GenModel3(eg.Out)
		}
	}

	for _, kp := range s.KafkaProducers {
		appDeps = append(appDeps, KafkaDep{Name: kp.Name, Kind: "producer"})
	}

	for _, kc := range s.KafkaConsumers {
		appDeps = append(appDeps, KafkaDep{Name: kc.Name, Kind: "consumer"})
	}

	/*
		for _, c := range s.Clients {
			clientDeps = append(clientDeps, c)
		}
	*/

	for _, e := range s.Endpoints {
		g.GenModel3(e.In)
		g.GenModel3(e.Out)
	}

	if s.Rpc == describer.GoNative {
		writeTemplate(fsys, service_net_rpc_tpl, opts, nil, tplGenPath[service_net_rpc_tpl])
		writeTemplate(fsys, app_tpl, opts, nil, tplGenPath[app_tpl])
		writeTemplate(fsys, client_net_rpc, opts, nil, tplGenPath[client_net_rpc])
	} else if s.Rpc == describer.HTTP {
		writeTemplate(fsys, service_http_tpl, opts, nil, tplGenPath[service_http_tpl])
		writeTemplate(fsys, app_tpl, opts, nil, tplGenPath[app_tpl])
		writeTemplate(fsys, client_http_tpl, opts, nil, tplGenPath[client_http_tpl])
	}

	writeTemplate(fsys, service_config_tpl, opts, nil, tplGenPath[service_config_tpl])
	writeTemplate(fsys, app_config_tpl, opts, nil, tplGenPath[app_config_tpl])

	fmt.Println("Generated Go files for service", s.Name)

	return nil
}
