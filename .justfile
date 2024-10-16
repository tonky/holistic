r:
  go run main.go

rr:
  fd -E gen '(.go|.tpl)' | entr -cr go run main.go 

rsp:
  ls gen/services/pizzeria/http/*.go | entr -cr go run gen/services/pizzeria/http/server_http.go

ms:
  ls tst/server/*.go | entr -cr go run tst/server/pizzeria.go
