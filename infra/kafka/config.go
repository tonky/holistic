package kafka

type Config struct {
	Brokers []string `required:"true" split_words:"true"`
}
