package tele

type IMetrics interface {
	IncCounter(name string, value int, tags map[string]string)
	IncHistogram(name string, value float64, tags map[string]string)
}
