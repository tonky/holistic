r:
  go run main.go

rr:
  fd -E gen '(.go|.tpl)' | entr -cr go run main.go 
