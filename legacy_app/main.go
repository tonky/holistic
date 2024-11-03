package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tonky/holistic/apps/pizzeria"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/pizzeria"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/do/v2"
)

type handlers struct {
	app *pizzeria.App
}

func main() {
	deps := do.New()

	l := logger.Slog{}

	conf, err := pizzeria.NewEnvConfig()
	if err != nil {
		panic(err)
	}

	po, err := pizzeria.NewPostgresOrdererRepository(l, conf.PostgresOrderer)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(deps, &l)
	do.ProvideValue(deps, po)
	do.Provide(deps, pizzeria.NewMemoryOrderProducerRepository)

	pa, err := pizzeria.NewApp(deps)
	if err != nil {
		panic(err)
	}

	h := handlers{app: pa}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", h.get)
	r.Put("/orders/{orderID}", h.update)
	r.Get("/orders/{orderID}", h.getOrder)
	fmt.Println("Starting server on :3000...")

	http.ListenAndServe(":3000", r)
}

func (h handlers) getOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderID")

	id, err := food.NewOrderID(orderIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.app.ReadOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ob, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(ob)
}

func (h handlers) get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func (h handlers) update(w http.ResponseWriter, r *http.Request) {
	id, err := food.NewOrderID(chi.URLParam(r, "orderID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateBody svc.UpdateOrder
	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateBody.ID = id

	order, errAppUpdate := h.app.UpdateOrder(r.Context(), updateBody.ToApp())
	if errAppUpdate != nil {
		http.Error(w, errAppUpdate.Error(), http.StatusBadRequest)
		return
	}

	ob, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(ob)
}
