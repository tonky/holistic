package kafka

type Config struct {
	Brokers []string `required:"true" split_words:"true" default:"localhost:19092"`
}

var liveKafkaBrokers = []string{"cloud-server-1:19092", "cloud-server-2:19092"}
var stagingKafkaBrokers = []string{"cloud-server-staging-1:19092"}

func EnvConfig(env string) Config {
	brokers := []string{"localhost:19092"}

	if env == "live" {
		brokers = liveKafkaBrokers
	} else if env == "staging" {
		brokers = stagingKafkaBrokers
	}

	return Config{
		Brokers: brokers,
	}
}
