package clients

func ConfigForEnv(service_name, env string) Config {
	host := "localhost"

	if service_name == "pricing" {
		host = "http://localhost"
	}

	if env == "local" || env == "dev" || env == "test" {
		return Config{
			Host: host,
			Port: servicePortsLocal[service_name],
		}
	}

	return Config{
		Host: host,
		Port: servicePortsLocal[service_name],
	}
}

var servicePortsLocal = map[string]int{
	"pizzeria":   1234,
	"pricing":    1235,
	"accounting": 1236,
	"shipping":   1237,
}
