package generator

import (
	"fmt"
	"tonky/holistic/generator/services"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func GenKafka(tds []services.TopicDesc) {
	template_dir := "templates"
	kafka_producer_tpl := "kafka_producer_v2.tpl"

	fsys := scriggo.Files{
		kafka_producer_tpl: readContent(template_dir, kafka_producer_tpl),
	}

	for _, t := range tds {
		fmt.Printf("Generating files for kafka topic: %s\n", t.TopicName)

		// tplGenPath := map[string]string{ kafka_producer_tpl: fmt.Sprintf("infra/kafkaProducer/%s.go", t.SnakeFileName()), }

		opts := &scriggo.BuildOptions{
			Globals: native.Declarations{
				"mod": "tonky/holistic",
				"cap": builtin.Capitalize,
				"kp":  &t,
			},
		}

		outFile := fmt.Sprintf("infra/kafkaProducer/%s.go", toSnakeCase(t.Name))
		writeTemplate(fsys, kafka_producer_tpl, opts, nil, outFile)

		fmt.Println("Generated kafka files")
	}
}

type KafkaDep struct {
	Name string
	Kind string
}

func (kd KafkaDep) InterfaceName() string {
	packageName := "kafkaConsumer"
	if kd.Kind == "producer" {
		packageName = "kafkaProducer"
	}

	return packageName + "." + "I" + builtin.Capitalize(kd.Name)
}

func (kd KafkaDep) AppVarName() string {
	return builtin.Capitalize(kd.Name) + builtin.Capitalize(kd.Kind)
}
