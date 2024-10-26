package services

import "strings"

type Postgres struct {
	Name      string
	Interface []Interface
}

func (p Postgres) StructName() string {
	return "Postgres" + strings.ToUpper(p.Name[0:1]) + p.Name[1:]
}

func (p Postgres) AppVarName() string {
	return p.Name + "Repo"
}

func (p Postgres) InterfaceName() string {
	return strings.ToUpper(p.Name[0:1]) + p.Name[1:] + "Repository"
}

type KafkaProducer struct {
	Name  string
	Topic string
	Model string
}

func (kp KafkaProducer) StructName() string {
	return "Kafka" + kp.InterfaceName()
}

func (kp KafkaProducer) AppVarName() string {
	return kp.Name + "Producer"
}

func (kp KafkaProducer) InterfaceName() string {
	return strings.ToUpper(kp.Name[0:1]) + kp.Name[1:] + "Producer"
}
