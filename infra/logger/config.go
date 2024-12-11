package logger

type LogFormat string

const (
	Text LogFormat = "text"
	JSON LogFormat = "json"
)

type Config struct {
	Level     string
	LogFormat LogFormat `default:"text"`
}
