r:
  go run main.go

rr:
  fd -E apps -E services -E clients -E domain -E infra '(.go|.tpl)' | entr -cr go run main.go 

t:
  PIZZERIA_PORT=1246 go test ./tests

ts:
  ls services/pizzeria/*.go | entr -cr go run tst/server/pizzeria.go

tc:
  PIZZERIA_PORT=1245 go run tst/main.go

ms:
  PIZZERIA_SHOULD_MOCK_APP=true just ts

pgc:
  psql postgres://postgres:postgres@localhost

pgs:
  docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

rpw:
  docker exec -it redpanda-0 rpk topic consume food.order.created -o -1

stats:
  wc -l apps/accounting/gen*.go services/accounting/* infra/postgres/* infra/kafkaConsumer/consumer.go infra/kafkaConsumer/food_order_shipped_consumer.go infra/kafkaProducer/producer.go infra/kafkaProducer/accounting_order_paid_producer.go clients/config.go clients/pricingClient/gen_client_http.go 

  wc -l apps/accounting/{app,orders_repository}.go
