#  Service generator
![Overview](https://github.com/tonky/holistic/blob/main/docs/arch_v2.svg)

## Stats
Lines of code in a production service
 - Total:          43355
    - Infra & domain: 27694(64%)
    - Application:    15661(36%)

Number of commits
 - Total:          551
   - Infra & domain: 394(70%)
   - Application:    168(30%)

```
$ git log --oneline -- ./internal/application | wc -l 
168

$ find ./internal/application  -name '*.go' | xargs  wc -l
15661 total

git log --oneline -- ./internal/domain ./internal/infra/{clients,http,kafka,marshaling,observability,rabbitmq} ./ahoy ./.github | wc -l
394

$ find ./internal/domain ./internal/infra/{clients,http,kafka,marshaling,observability,rabbitmq} -name '*.go' | xargs  wc -l
27694 total
```

## Declare service
```
func AccountingService() describer.Service {
	getOrder := describer.Endpoint{
		Name:   "order",
		Method: describer.Read,
		In:     describer.Inputs{Name: "food.OrderID"},
		Out: map[describer.ResponseType]describer.ResponseObject{
			describer.ResponseOK: "accounting.Order",
		},
	}

	return describer.Service{
		Name:           "accounting",
		Rpc:            describer.GoNative,
		Dependencies:   describer.Struct,
		ConfigItems:    []describer.ConfigItem{{Name: "KafkaConsumptionRPS", Typ: "int"}},
		Endpoints:      []describer.Endpoint{getOrder},
		KafkaConsumers: []describer.TopicDesc{TopicFoodOrderUpdated},
		KafkaProducers: []describer.TopicDesc{TopicAccountingOrderPaid},
		Postgres: []describer.Postgres{{
			Name: "orderer",
			Methods: []describer.InterfaceMethod{
				{
					Name: "SaveFinishedOrder",
					Arg:  describer.InfraObject{Name: "orderID", Typ: "accounting.Order"},
					Ret:  describer.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
				{
					Name: "ReadOrderByFoodID",
					Arg:  describer.InfraObject{Name: "newOrder", Typ: "food.OrderID"},
					Ret:  describer.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
			},
		}},
		Clients: []describer.Client{
			{
				VarName: "pricingClient",
				IName:   "IPricingClient",
			},
		},
	}
}
```

## How application code looks
```
// apps/accounting/app.go

// Kafka 'food.order.updated' consumer handler
func (a *App) FoodOrderUpdatedProcessor(ctx context.Context, in food.Order) error {
	a.Logger.Info("AccountingApp.FoodOrderUpdatedProcessor got: ", in)

	if !in.IsFinal {
		return nil
	}

	orderPrice, err := a.Clients.PricingClient.ReadOrder(ctx, in.ID)
	if err != nil {
		return err
	}

	paidOrder := accounting.Order{
		ID:   in.ID,
		Cost: orderPrice.Cost,
	}

	_, errSave := a.Deps.OrdererRepo.SaveFinishedOrder(ctx, paidOrder)
	if errSave != nil {
		return errSave
	}

	return a.Deps.AccountingOrderPaidProducer.ProduceAccountingOrderPaid(ctx, paidOrder)
}
```

## Overview of this repo
```
├── apps # generated application code
│   ├── accounting
│   │   ├── app.go # <-- app logic goes here
│   │   ├── gen_accounting_app.go
│   │   ├── gen_config.go
│   │   ├── gen_orderer_repository_postgres.go # 'New' and interface
│   │   ├── orders_memory_repo.go # implemented by developer
│   │   └── orders_repository.go # implemented by developer
│   ├── pizzeria
│   ├── pricing
│   └── shipping
├── clients # generated clients for services
│   ├── accountingClient
│   │   ├── gen_client_net_rpc.go # based on service transport
│   │   └── mock.go
│   ├── config.go
│   ├── from_env.go
│   ├── pizzeriaClient
│   ├── pricingClient
│   │   └── gen_client_http.go
│   └── shippingClient
├── decl # declarations of domain types, kafka and services
│   ├── accounting.go
│   ├── domain_types.go # <-- domain types declarations
│   ├── kafka.go # <-- available kafka topics
│   ├── pizzeria.go
│   ├── pricing.go
│   └── shipping.go
├── domain # generated domain code
│   ├── accounting
│   │   ├── order.go
│   │   └── orderprice.go
│   ├── food
│   ├── pricing
│   └── shipping
├── infra
│   ├── kafka
│   │   ├── config.go # common config base for consumers/producers
│   ├── kafkaConsumer # generated consumers per topic
│   │   ├── accounting_order_paid_consumer.go
│   │   ├── consumer.go
│   │   ├── food_order_created_consumer.go
│   │   ├── food_order_paid_consumer.go
│   │   ├── food_order_shipped_consumer.go
│   │   ├── food_order_updated_consumer.go
│   │   └── shipping_order_shipped_consumer.go
│   ├── kafkaProducer # generated producers per topic
│   │   ├── accounting_order_paid_producer.go
│   │   ├── food_order_created_producer.go
│   │   ├── food_order_paid_producer.go
│   │   ├── food_order_shipped_producer.go
│   │   ├── food_order_updated_producer.go
│   │   ├── producer.go
│   │   └── shipping_order_shipped_producer.go
│   ├── logger
│   │   └── slog.go
│   ├── postgres
│   │   └── pgx.go
├── services
│   ├── accounting
│   ├── pizzeria
│   │   ├── config.go
│   │   ├── new_order.go # custom 'in' type by developer
│   │   └── server_net_rpc.go # transport-specific server
│   ├── pricing
│   └── shipping
│       ├── config.go
│       └── server_http.go
├── tests
│   ├── integration_test.go <-- see example of 4 services test
│   └── order_e2e_test.go
```

## Scope of work for PoC
- [x] Domain types
- [x] Server generator
- [x] Server: multiple transports
- [x] Server: opionated DTO
- [x] Server: custom config variables
- [x] Client generator
- [x] Client: multiple transports
- [x] Kafka: topics description
- [x] Kafka: consumers generation
- [x] Kafka: producers generation

- [x] App: custom config variables
- [x] App dependency: Postgres
- [x] App dependency: Kafka consumer
- [x] App dependency: Kafka producer
- [x] Integration test with postgres and kafka

## Plans
- [ ] Path to peacemeal conversion of existing services
- [ ] Versioning: domain, services
- [ ] Kafka: integrate Schema Registry
- [ ] Service: discovery
- [ ] Service: deployment-wide configuration
- [ ] Service: API docs generator based on transport
- [ ] Service: Cron jobs
- [ ] Service: Background jobs
- [ ] Service: ACL
- [ ] Service: Traces, logs and metrics
- [ ] Service: Configurable serialization
- [ ] App: Traces, logs and metrics
- [ ] App: Middlewares
- [ ] Caching
- [ ] CLI tools for better DX
