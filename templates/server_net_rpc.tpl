// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name}}

import (
    "context"
    "fmt"
	"log"
	"net"
	"net/rpc"

	{% if service.KafkaProducers %}
	"tonky/holistic/infra/kafkaProducer"
	{% end %}
	{% if service.KafkaConsumers %}
	"tonky/holistic/infra/kafkaConsumer"
	{% end %}

    // "github.com/go-playground/validator/v10"
    {% if service.Dependencies == "samber_do" %}
    "github.com/samber/do/v2"
    {% end %}

    {% for id in app_deps %}
        {% if id.AppImportPackageName() == "app" %}
        {% else if id.PackageName() == "local" %} 
        {% else if id.PackageName() == "kafkaProducer" %} 
        {% else if id.PackageName() == "kafkaConsumer" %} 
        {% else %}
	"tonky/holistic/infra/{{ id.PackageName() }}"
        {% end %}
    {% end %}

	{% for imp in service.ClientImports() %}
    {% if imp.Alias == "svc" %}
        {% continue %}
    {% end if %}
	{% if imp.Alias %}{{ imp.Alias }} {% end if %}"tonky/holistic/{{ imp.RelPath }}"
	{% end %}
	app "tonky/holistic/apps/{{ service.Name }}"
	"tonky/holistic/infra/logger"
    {% for d in client_deps.Dedup() %}
	"tonky/holistic/clients/{{ d.AppVarName() }}"
    {% end %}
)

type {{ cap(service.Name) }} struct {
    config Config
    {% if service.Dependencies == "samber_do" %}
    deps do.Injector
    {% else if service.Dependencies == "plain_struct" %}
    deps app.Deps
    {% end %}
    app app.App
}

{% for h in handlers %}
func (h {{ cap(service.Name) }}) {{h.FuncName()}}(arg {{ h.In }}, reply *{{ h.Out.ok }}) error {
    res, err := h.app.{{h.FuncName()}}(context.TODO(), arg{{ h.In.SvcToApp() }})
    if err != nil {
        return err
    }

    *reply = res

    return nil
}

{% end %}

{% if service.Dependencies == "plain_struct" %}
    {% if client_deps %}
func New(deps app.Deps, clients app.Clients) (ServiceStarter, error) {
    {% else %}
func New(deps app.Deps) (ServiceStarter, error) {
    {% end %}
{% else if service.Dependencies == "samber_do" %}
func New(deps app.Deps) (ServiceStarter, error) {
{% end %}
	cfg, err := NewEnvConfig()
    if err != nil {
        return nil, err
    }

{% if service.Dependencies == "plain_struct" %}
    {% if client_deps %}
    application, appErr := app.NewApp(deps, clients)
    {% else %}
    application, appErr := app.NewApp(deps)
    {% end %}
{% else if service.Dependencies == "samber_do" %}
    application, appErr := app.NewApp(deps)
{% end %}
    if appErr != nil {
        return nil, appErr
    }

    handlers := {{ cap(service.Name) }}{deps: deps, config: cfg, app: *application}

    return handlers, nil
}

func (h {{ cap(service.Name) }}) Start() error {
	port := h.config.Port

    fmt.Printf(">> {{ service.Name }}.Start() config: %+v\n", h.config)
	{% if service.KafkaConsumers %}
    h.app.RunConsumers()
    {% end %}

	server := rpc.NewServer()
	server.Register(h)

	fmt.Println(">> starting server on port ", port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatal("listen error:", err)
    }

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			server.ServeConn(conn)
		}()
	}
}
// TODO: REMOVE
{% if service.Dependencies == "samber_do" %}
{% else if service.Dependencies == "plain_struct" %}
{% end %}

func NewFromEnv() (ServiceStarter, error) {
	cfg, err := NewEnvConfig()
    if err != nil {
        return nil, err
    }

{% if service.Dependencies == "samber_do" %}
    deps := do.New()

	do.ProvideValue(deps, cfg)

    l := logger.Slog{}
	do.ProvideValue(deps, &l)

    {% for ad in app_deps %}
    {% if ad.PackageName() == "local" %}
	{{ ad.AppVarName() }} := app.New{{ ad.StructName() }}(l)
    {% else %}
	{{ ad.AppVarName() }}, err := {% if ad.PackageName() != "local" %}{{ ad.AppImportPackageName() }}.{% else %}app.{% end %}New{{ ad.StructName() }}(l, cfg.App.{{ ad.ConfigVarName() }})
    if err != nil {
        return nil, err
    }
    {% end %}

	do.ProvideValue(deps, {{ ad.AppVarName() }})

    {% end %}
    {% for d in client_deps.Dedup() %}
	do.ProvideValue(deps, {{ d.AppVarName() }}.NewFromEnv(cfg.Environment))
    {% end %}
{% else if service.Dependencies == "plain_struct" %}
    // plain deps code here
    deps := app.Deps{
        Logger: &logger.Slog{},
    }

    {% for ad in app_deps %}
        {% if ad.PackageName() == "local" %}
	deps.{{ cap(ad.AppVarName()) }} = app.New{{ ad.StructName() }}(l)
        {% else %}
	{{ ad.AppVarName() }}, err := {% if ad.PackageName() != "local" %}{{ ad.AppImportPackageName() }}.{% else %}app.{% end %}New{{ ad.StructName() }}(*deps.Logger, cfg.App.{{ ad.ConfigVarName() }})
    if err != nil {
        return nil, err
    }

	deps.{{ cap(ad.AppVarName()) }} = {{ ad.AppVarName() }}
        {% end %}

    {% end for %}
    {% if client_deps %}
    clients := app.Clients {
        {% for d in client_deps.Dedup() %}
        {{ cap(d.AppVarName()) }}: {{ d.AppVarName() }}.NewFromEnv(cfg.Environment),
        {% end %}
    }
    {% end %}
{% end %}

    {% if client_deps %}
    application, appErr := app.NewApp(deps, clients)
    {% else %}
    application, appErr := app.NewApp(deps)
    {% end %}
    if appErr != nil {
        return nil, appErr
    }

    handlers := {{ cap(service.Name) }}{deps: deps, config: cfg, app: *application}

    return handlers, nil
}

func (h {{ cap(service.Name) }}) Config() Config {
    return h.config
}

func (h {{ cap(service.Name) }}) Deps() app.Deps {
    return h.deps
}


type ServiceStarter interface {
    Start() error
    Config() Config
    Deps() app.Deps
}
