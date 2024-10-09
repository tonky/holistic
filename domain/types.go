package domain

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
