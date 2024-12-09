// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name}}

import (
    "context"
    "fmt"
	"net/http"
    "encoding/json"

	"github.com/go-chi/chi/v5"

	app "{{ modulePath }}/apps/{{ service.Name }}"

    {% for ci in service.AbsImports(ctx) %}
	"{{ ci }}"
    {% end %}
)

type Deps struct {
    Logger {{ service.Logger.Interface.GoQualifiedModel() }}
}

type handlers struct {
    config Config
    deps Deps
    app app.App
}

{% for h in service.Endpoints %}
func (h handlers) {{h.Name}}() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in {{ h.In.GoQualifiedModel() }}
        var out {{ h.Out.GoQualifiedModel() }}

        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        out, appErr := h.app.{{h.Name}}(context.TODO(), in)

        if appErr != nil {
            http.Error(w, appErr.Error(), http.StatusBadRequest)
            return
        }

        byteResp, err := json.Marshal(out)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Write(byteResp)
    }
}

{% end %}

func NewFromEnv() (ServiceStarter, error) {
	cfg := MustEnvConfig()

    log := {{ service.Logger.Model.Package() }}.Default().With("service", "{{ service.Name }}")

    log.Debug("NewFromEnv()", "config", cfg)

    deps, err := app.DepsFromConf(cfg.App)
    if err != nil {
        return nil, err
    }

    deps.Logger = log

    application, appErr := app.NewApp(deps)
    if appErr != nil {
        return nil, appErr
    }

    handlers := handlers{config: cfg, app: *application, deps: Deps{Logger: log}}

    return handlers, nil
}

func NewWithAppDeps(serviceDeps Deps, appDeps app.Deps) (ServiceStarter, error) {
	cfg := MustEnvConfig()

    serviceDeps.Logger.Debug("{{ service.Name }}.NewWithAppDeps()", "config", cfg)

    appDeps.Logger = serviceDeps.Logger

    application, appErr := app.NewApp(appDeps)
    if appErr != nil {
        return nil, appErr
    }

    return handlers{config: cfg, app: *application, deps: serviceDeps}, nil
}

func DepsFromConf(conf Config) Deps {
    return Deps {
        Logger: {{ service.Logger.Model.Package() }}.NewFromConfig(conf.Logger),
    }
}


func (h handlers) Start() error {
    h.deps.Logger.Info("Start()", "config", h.config)
	{% if service.KafkaConsumers %}
    h.app.RunConsumers()
    {% end %}

	r := chi.NewRouter()

    {% for h in service.Endpoints %}
	r.Post("/{{ cap(h.Name) }}", h.{{ h.Name }}())
    {% end %}

    return http.ListenAndServe(fmt.Sprintf(":%d", h.config.Port), r)
}

func (h handlers) Config() Config {
    return h.config
}

func (h handlers) App() app.App {
    return h.app
}

func (h handlers) Deps() Deps {
    return h.deps
}

type ServiceStarter interface {
    Start() error
    Config() Config
    App() app.App
    Deps() Deps
}
