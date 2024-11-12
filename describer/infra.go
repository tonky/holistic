package describer

import "strings"

type Postgres struct {
	Name    string
	Methods []InterfaceMethod
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

func (p Postgres) PackageName() string {
	return "postgres"
}

func (p Postgres) AppImportPackageName() string {
	return "app"
}

func (p Postgres) ConfigVarName() string {
	return p.StructName()
}

func (p Postgres) ConfigVarType() string {
	return "postgres.Config"
}

func (p Postgres) Imports(s Service) []ClientImport {
	all := []ClientImport{}
	for _, m := range p.Methods {
		all = append(all, m.Imports(s)...)
	}

	var out []ClientImport
	seen := map[string]bool{}

	for _, dep := range all {
		if seen[dep.RelPath] {
			continue
		}

		seen[dep.RelPath] = true

		out = append(out, dep)
	}

	return out
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
