package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // config := 
    handlers := handlers{deps: dependencies}
    
    {% for h in handlers %}
    r.{{h.Method.HttpName()}}("{{ h.In.Path() }}", handlers.{{h.FuncName()}})
    {% end %}

    http.ListenAndServe(":{{ port }}", r)
}

{% for h in handlers %}
type {{h.FuncName()}}Request struct {
    {% for input in h.In %}
    {{input.What.FieldName()}} {{input.What.Name}} {% if input.Validation %}`validate:"{{ input.Validation }}"`{% end %}
    {% end %}
} 

validate.RegisterStructValidation({{h.FuncName()}}RequestValidator, {{h.FuncName()}}Request{})

func {{h.FuncName()}}RequestValidator()(sl validator.StructLevel) {
    // add custom validator here, if specified on endpoint
    return nil
}

type {{h.FuncName()}}Response struct {} 

// swagger definitions here
func (h handlers) {{h.FuncName()}}(w http.ResponseWriter, r *http.Request)
    func (in {{h.FuncName()}}Request) (responseModel, error) {
        requestModel := {{h.FuncName()}}RequestParse(r)

        if err := validate.Struct(requestModel); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write(w, err)
            
            return
        }

        responseModel, err := app.{{h.FuncName()}}(r.context, requestModel, h.deps)

        response, err := json.Marshal(responseModel)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(response)
    }
}

func {{h.FuncName()}}RequestParse(r *http.Request) {{h.FuncName()}}Request {
    var model {{h.FuncName()}}Request{}

    {% for param in h.In %}
        {% if param.Where == "path" %}
    model.{{ param.URLParamName() }} = chi.URLParam(r, "{{ param.URLParamName() }}")
        {% end %}
    {% end %}

    return model
}

{% end %}

type handlers struct {
    deps *do.Injector
}
