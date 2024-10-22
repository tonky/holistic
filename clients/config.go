package clients

import "fmt"

type Config struct {
	Host string
	Port int
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
