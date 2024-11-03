// AUTOGENERATED! DO NOT EDIT.
package legacy

import (
    "context"
	"net/http"
    "encoding/json"

	"github.com/go-chi/chi/v5"

	app "tonky/holistic/apps/legacy"
)

type handlers struct {
    config Config
    app app.App
}

func (h handlers) ReadOrder() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in OrderID

        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        appArg, err := in.ToApp()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        respModel, err := h.app.ReadOrder(context.TODO(), appArg)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        byteResp, err := json.Marshal(respModel)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Write(byteResp)
    }
}

func (h handlers) CreateOrder() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in NewOrder

        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        appArg, err := in.ToApp()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        respModel, err := h.app.CreateOrder(context.TODO(), appArg)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        byteResp, err := json.Marshal(respModel)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Write(byteResp)
    }
}

func (h handlers) UpdateOrder() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in UpdateOrder

        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        appArg, err := in.ToApp()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        respModel, err := h.app.UpdateOrder(context.TODO(), appArg)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        byteResp, err := json.Marshal(respModel)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Write(byteResp)
    }
}


func NewLegacy(conf Config, appDeps app.Deps) (*chi.Mux, error) {
	a, err := app.NewApp(appDeps)
	if err != nil {
		panic(err)
	}

	h := handlers{app: a, config: conf}

	r := chi.NewRouter()

	r.Post("/GetOrder", h.ReadOrder())
	r.Post("/PostOrder", h.CreateOrder())
	r.Post("/PutOrder", h.UpdateOrder())

    return r, nil
}
