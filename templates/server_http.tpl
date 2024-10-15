package main

import (
	"encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/go-playground/validator/v10"
    "github.com/samber/do/v2"

	"tonky/holistic/gen/domain/food"
	"tonky/holistic/gen/services/{{ service_name }}/app"
)

type Config struct {
	Port int
}

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    dependencies := do.New()
    do.ProvideValue(dependencies, &Config{
        Port: 4242,
    })

    // car, err := do.Invoke[*Car](i)
    // if err != nil { log.Fatal(err.Error()) }
    // car.Start()

    handlers := handlers{deps: dependencies}
    
    {% for h in handlers %}
    r.{{h.Method.HttpName()}}("{{ h.In }}", handlers.{{h.FuncName()}})
    {% end %}

    http.ListenAndServe(":{{ port }}", r)
}

{% for h in handlers %}
type {{h.FuncName()}}Request struct {
    {{ h.In.ModelName() }} {{ h.In.Name }}
} 

func {{h.FuncName()}}RequestValidator()(sl validator.StructLevel) {
    // add custom validator here, if specified on endpoint
    return nil
}

type {{h.FuncName()}}Response struct {} 

// swagger definitions here
func (h handlers) {{h.FuncName()}}(w http.ResponseWriter, r *http.Request) {
    {#
    requestModel := {{h.FuncName()}}RequestParse(r)

	validate := validator.New(validator.WithRequiredStructEnabled())
    if err := validate.Struct(requestModel); err != nil {
        w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
        
        return
    }
    #}
    application := app.New(h.deps)

    {# appArgs := dtoToAppArgs{{h.FuncName()}}(requestModel) #}
    appArgs := {{ h.In }}{}

    responseModel, err := application.{{h.FuncName()}}(r.Context(), appArgs)

    response, err := json.Marshal(responseModel)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(response)
}

{#
func {{h.FuncName()}}RequestParse(r *http.Request) {{h.FuncName()}}Request {
    model := {{h.FuncName()}}Request{}

    {% for param in h.In %}
        {% if param.Where == "path" %}
    {{ param.URLParamName() }}Str := chi.URLParam(r, "{{ param.URLParamName() }}")
    {{ param.URLParamName() }}, err := {{param.What.Name}}New({{ param.URLParamName() }}Str)
    if err != nil {
        return model
    }

    model.{{ param.URLParamName() }} = {{ param.URLParamName() }}
        {% end %}
    {% end %}

    return model
}
#}
func dtoToAppArgs{{h.FuncName()}}(requestModel {{h.FuncName()}}Request) {
    // return resp
}

{% end %}

type handlers struct {
    deps do.Injector
}
