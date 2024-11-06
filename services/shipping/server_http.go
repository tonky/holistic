package shipping

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	app "tonky/holistic/apps/shipping"

	"tonky/holistic/domain/food"
	"tonky/holistic/domain/shipping"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/logger"
)

type handlers struct {
	config Config
	app    app.App
	deps   app.Deps
	mux    *chi.Mux
}

func (h handlers) ReadOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in food.OrderID
		var out shipping.Order

		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var appErr error

		out, appErr = h.app.ReadOrder(context.TODO(), in)
		if appErr != nil {
			http.Error(w, appErr.Error(), http.StatusInternalServerError)
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

func NewMux(conf Config, deps app.Deps) (*chi.Mux, error) {
	a, err := app.NewApp(deps)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	h := handlers{deps: deps, app: a, config: conf, mux: r}

	r.Post("/GetOrder", h.ReadOrder())

	return r, nil
}

func NewFromEnv() (ServiceStarter, error) {
	cfg, err := NewEnvConfig()
	if err != nil {
		return nil, err
	}

	deps := app.Deps{
		Logger: &logger.Slog{},
		Config: cfg.App,
	}

	ordererRepo, err := app.NewPostgresOrderer(*deps.Logger, cfg.App.PostgresOrderer)
	if err != nil {
		return nil, err
	}
	deps.OrdererRepo = ordererRepo
	ShippingOrderShippedProducer, err := kafkaProducer.NewShippingOrderShippedProducer(*deps.Logger, cfg.App.Kafka)
	if err != nil {
		return nil, err
	}
	deps.ShippingOrderShippedProducer = ShippingOrderShippedProducer
	AccountingOrderPaidConsumer, err := kafkaConsumer.NewAccountingOrderPaidConsumer(*deps.Logger, cfg.App.Kafka)
	if err != nil {
		return nil, err
	}
	deps.AccountingOrderPaidConsumer = AccountingOrderPaidConsumer

	application, appErr := app.NewApp(deps)
	if appErr != nil {
		return nil, appErr
	}

	mux, err := NewMux(cfg, deps)
	if err != nil {
		return nil, err
	}

	handlers := handlers{config: cfg, app: application, deps: deps, mux: mux}

	return handlers, nil
}

func (h handlers) Config() Config {
	return h.config
}

func (h handlers) Start() error {
	fmt.Printf(">> shipping.Start() config: %+v\n", h.config)

	return http.ListenAndServe(fmt.Sprintf(":%d", h.config.Port), h.mux)
}

func (h handlers) Deps() app.Deps {
	return h.deps
}

func (h handlers) Mux() *chi.Mux {
	return h.mux
}

type ServiceStarter interface {
	Start() error
	Config() Config
	Deps() app.Deps
	Mux() *chi.Mux
}
