package kafkaProducer

type Config struct {
	Brokers []string `required:"true" split_words:"true"`
}
