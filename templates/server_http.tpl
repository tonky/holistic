package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    {% for h in handlers %}

    // swagger definitions here
    r.{{ h.Method.HttpName() }}("{{ h.In.Path() }}", func(w http.ResponseWriter, r *http.Request) {
        // parse request into domain types
        // run validations, if defined

        // call http handler
        responseModel, err := {{h.FuncName()}}({{h.FuncArgs()}})

        response, err := json.Marshal(responseModel)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("{{ h.Out["server_error"] }}"))
        }

        json.NewEncoder(w).Encode(response)
    })
    {% end %}

    http.ListenAndServe(":{{ port }}", r)
}

{% for h in handlers %}

func {{h.FuncName()}}({{h.FuncArgs()}}) ({{h.Out["ok"]}}, error) {
    var responseModel {{ h.Out["ok"] }}

    return responseModel
}
{% end %}