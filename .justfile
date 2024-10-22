r:
  go run main.go

rr:
  fd -E services -E clients -E domain '(.go|.tpl)' | entr -cr go run main.go 

t:
  go test tests/*

ts:
  ls services/pizzeria/http/*.go | entr -cr go run tst/server/pizzeria.go

ms:
  PIZZERIA_SHOULD_MOCK_APP=true just ts

pgc:
  psql postgres://postgres:postgres@localhost

pgs:
  docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
