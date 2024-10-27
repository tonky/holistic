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

func (k KafkaProducer) Imports() []ClientImport {
	return []ClientImport{importFromModel(k.Model)}
}

type KafkaConsumer struct {
	Name  string
	Topic string
	Model string
}

func (k KafkaConsumer) StructName() string {
	return "Kafka" + k.InterfaceName()
}

func (k KafkaConsumer) AppVarName() string {
	return k.Name + "Consumer"
}

func (k KafkaConsumer) InterfaceName() string {
	return strings.ToUpper(k.Name[0:1]) + k.Name[1:] + "Consumer"
}

func (k KafkaConsumer) Imports() []ClientImport {
	return []ClientImport{importFromModel(k.Model)}
}

func importFromModel(model string) ClientImport {
	split := strings.Split(model, ".")

	return ClientImport{RelPath: "domain/" + split[0]}
}
