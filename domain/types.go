package domain

import "github.com/open2b/scriggo/builtin"

type Kind int

const (
	Int Kind = iota
	Float
	String
	UUID
	StringList
	ObjectList
)

type FieldName string

type Object struct {
	Name   string
	Fields map[FieldName]Kind
}

func (o Object) FieldName() string {
	tmp := builtin.Split(o.Name, ".")

	return tmp[len(tmp)-1]
}
