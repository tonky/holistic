package generator

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"tonky/holistic/describer"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func GenService(s describer.Service) {
	fmt.Printf("Generating %s service Go files\n", s.Name)

	template_dir := "templates"
	service_net_rpc_tpl := "server_net_rpc.tpl"
	service_http_tpl := "server_http.tpl"
	client_net_rpc := "client_net_rpc.tpl"
	client_http_tpl := "client_http.tpl"
	service_config_tpl := "service_config.tpl"
	app_config_tpl := "app_config.tpl"
	app_tpl := "app.tpl"
	app_plain_tpl := "app_plain.tpl"
	repo_pg_tpl := "repository_postgres.tpl"
	repo_generic_tpl := "repository_generic.tpl"
	// kafka_consumer_tpl := "kafka_consumer.tpl"

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
		app_plain_tpl:       fmt.Sprintf("apps/%s/gen_%s_app.go", s.Name, s.Name),
		app_config_tpl:      fmt.Sprintf("apps/%s/gen_config.go", s.Name),
	}

	fsys := scriggo.Files{
		service_net_rpc_tpl: readContent(template_dir, service_net_rpc_tpl),
		service_http_tpl:    readContent(template_dir, service_http_tpl),
		client_net_rpc:      readContent(template_dir, client_net_rpc),
		client_http_tpl:     readContent(template_dir, client_http_tpl),
		service_config_tpl:  readContent(template_dir, service_config_tpl),
		app_tpl:             readContent(template_dir, app_tpl),
		app_plain_tpl:       readContent(template_dir, app_plain_tpl),
		app_config_tpl:      readContent(template_dir, app_config_tpl),
		repo_pg_tpl:         readContent(template_dir, repo_pg_tpl),
		// kafka_producer_tpl:  readContent(template_dir, kafka_producer_tpl),
		// kafka_consumer_tpl:  readContent(template_dir, kafka_consumer_tpl),
	}

	opts := &scriggo.BuildOptions{
		Globals: native.Declarations{
			"mod":              "tonky/holistic",
			"service":          &s,
			"cap":              builtin.Capitalize,
			"port":             (*int)(nil),
			"handlers":         &s.Endpoints,
			"config_items":     &s.ConfigItems,
			"app_config_items": &s.AppConfigItems,
			"infra":            &s.Infra,
		},
	}

	appDeps := AppDeps{}
	infraDeps := AppDeps{}
	interfaces := AppDeps{}
	clientDeps := AppDeps{}
	configImports := s.ClientImports()

	opts.Globals["app_deps"] = &appDeps
	opts.Globals["infra_deps"] = &infraDeps
	opts.Globals["interfaces"] = &interfaces
	opts.Globals["client_deps"] = &clientDeps
	opts.Globals["client_relative_imports"] = &configImports

	for _, pg := range s.Postgres {
		appDeps = append(appDeps, pg)
		infraDeps = append(infraDeps, pg)

		opts.Globals["repo"] = &pg
		// opts.Globals["kind"] = "postgres"

		outFile := fmt.Sprintf("apps/%s/gen_%s_repository_postgres.go", s.Name, pg.Name)
		writeTemplate(fsys, repo_pg_tpl, opts, nil, outFile)
	}

	for _, i := range s.Interfaces {
		opts.Globals["repo"] = &i
		appDeps = append(appDeps, &i)

		outFile := fmt.Sprintf("apps/%s/gen_%s.go", s.Name, toSnakeCase(i.Struct))
		writeTemplate(fsys, repo_generic_tpl, opts, nil, outFile)
	}

	for _, kp := range s.KafkaProducers {
		appDeps = append(appDeps, KafkaDep{Name: kp.Name, Kind: "producer"})
	}

	for _, kc := range s.KafkaConsumers {
		appDeps = append(appDeps, KafkaDep{Name: kc.Name, Kind: "consumer"})
	}

	for _, c := range s.Clients {
		clientDeps = append(clientDeps, c)
	}

	if s.Rpc == describer.GoNative {
		writeTemplate(fsys, service_net_rpc_tpl, opts, nil, tplGenPath[service_net_rpc_tpl])
		writeTemplate(fsys, app_tpl, opts, nil, tplGenPath[app_tpl])
		writeTemplate(fsys, client_net_rpc, opts, nil, tplGenPath[client_net_rpc])
	} else if s.Rpc == describer.HTTP {
		writeTemplate(fsys, service_http_tpl, opts, nil, tplGenPath[service_http_tpl])
		writeTemplate(fsys, app_plain_tpl, opts, nil, tplGenPath[app_plain_tpl])
		writeTemplate(fsys, client_http_tpl, opts, nil, tplGenPath[client_http_tpl])
	}

	writeTemplate(fsys, service_config_tpl, opts, nil, tplGenPath[service_config_tpl])
	writeTemplate(fsys, app_config_tpl, opts, nil, tplGenPath[app_config_tpl])

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
	StructName() string
	PackageName() string
	AppImportPackageName() string
	ConfigVarName() string
	ConfigVarType() string
}

type AppDeps []AppDep

func (ads AppDeps) Dedup() []AppDep {
	keys := make(map[string]bool)
	list := []AppDep{}

	for _, entry := range ads {
		if _, value := keys[entry.ConfigVarName()]; !value {
			keys[entry.ConfigVarName()] = true
			list = append(list, entry)
		}
	}

	return list
}

type InfraDep struct {
	Typ  string
	Name string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
