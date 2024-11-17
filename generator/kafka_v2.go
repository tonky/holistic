package generator

import (
	"fmt"
	"tonky/holistic/describer"
	"tonky/holistic/typs"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/builtin"
	"github.com/open2b/scriggo/native"
)

func GenKafka2(template_dir string, tds []describer.TopicDesc2) {
	ctx := typs.Object3{}

	kafka_producer_tpl := "kafka_producer_v3.tpl"
	kafka_consumer_tpl := "kafka_consumer_v2.tpl"

	fsys := scriggo.Files{
		kafka_producer_tpl: readContent(template_dir, kafka_producer_tpl),
		kafka_consumer_tpl: readContent(template_dir, kafka_consumer_tpl),
	}

	mustCreateDirs([]string{"./infra/kafkaProducer", "./infra/kafkaConsumer"})

	for _, t := range tds {
		fmt.Printf("Generating files for kafka topic: %s\n", t.TopicName)

		// tplGenPath := map[string]string{ kafka_producer_tpl: fmt.Sprintf("infra/kafkaProducer/%s.go", t.SnakeFileName()), }

		opts := &scriggo.BuildOptions{
			Globals: native.Declarations{
				// "mod": "tonky/holistic",
				"cap": builtin.Capitalize,
				// "kp":  &t,
				"topic": &t,
				"ctx":   &ctx,
			},
		}

		outFile := fmt.Sprintf("infra/kafkaProducer/%s_producer.go", toSnakeCase(t.Name))
		writeTemplate(fsys, kafka_producer_tpl, opts, nil, outFile)

		outFileC := fmt.Sprintf("infra/kafkaConsumer/%s_consumer.go", toSnakeCase(t.Name))
		writeTemplate(fsys, kafka_consumer_tpl, opts, nil, outFileC)

		fmt.Println("Generated kafka files")
	}
}

func (g ServiceGen) GenKafka2(tds []describer.TopicDesc2) {
	GenKafka2(g.TemplatePath, tds)
}
