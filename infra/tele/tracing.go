package tele

type ITracing interface {
	StartSpan(operationName string)
}
