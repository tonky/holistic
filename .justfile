r:
  go run main.go

rr:
  fd -E apps -E services -E clients -E domain '(.go|.tpl)' | entr -cr go run main.go 

t:
  PIZZERIA_PORT=1246 go test tests/*

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

rpcl:
  docker exec -it redpanda-0 rpk topic consume pizzeria.order -o -1 -n 1

stats:
  wc -l apps/pizzeria/order* apps/pizzeria/pizzeria_app.go
  wc -l apps/pizzeria/gen* services/pizzeria/*.go clients/pizzeria*.go
