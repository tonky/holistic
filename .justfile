r:
  go run main.go

rr:
  fd -E gen '(.go|.tpl)' | entr -cr go run main.go 

ts:
  ls gen/services/pizzeria/http/*.go | entr -cr go run tst/server/pizzeria.go

ms:
  ls tst/server/*.go | entr -cr go run tst/server/pizzeria.go
