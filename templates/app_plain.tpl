// AUTOGENERATED! DO NOT EDIT.

package {{ service.Name }}

import (
	{% if service.KafkaConsumers %}
	"context"
	{% end %}
	{% if service.KafkaProducers %}
	"tonky/holistic/infra/kafkaProducer"
	{% end %}
	{% if service.KafkaConsumers %}
	"tonky/holistic/infra/kafkaConsumer"
	{% end %}
	{% for d in client_deps %}
	"tonky/holistic/clients/{{ d.AppVarName() }}"
	{% end %}
	"tonky/holistic/infra/logger"
)

type Deps struct {
	Config Config
	Logger logger.ILogger
{% for ad in app_deps %}
    {{ cap(ad.AppVarName()) }} {{ ad.InterfaceName() }}
{% end %}
}

{% if client_deps %}
type Clients struct {
{% for d in client_deps.Dedup() %}
    {{ cap(d.AppVarName()) }} {{ d.AppVarName() }}.{{ d.InterfaceName() }}
{% end %}
}
{% end %}

type App struct {
	Deps		Deps
{% if client_deps %}
	Clients		Clients
{% end %}
}

{% if client_deps %}
func NewApp(deps Deps, clients Clients) (App, error) {
{% else %}
func NewApp(deps Deps) (App, error) {
{% end %}
	app := App{
		Deps:       deps,
{% if client_deps %}
		Clients: 	clients,
{% end %}
	}

	return app, nil
}

{% if service.KafkaConsumers %}
func (a App) RunConsumers() {
	a.Deps.Logger.Info(">> {{ service.Name}}.App.RunConsumers()")

	ctx := context.Background()
	{% for consumer in service.KafkaConsumers %}

	go func() {
		for err := range a.Deps.{{ cap(consumer.Name) }}Consumer.Run(ctx, a.{{ cap(consumer.Name) }}Processor) {
			a.Deps.Logger.Warn(err.Error())
		}
	}()
	{% end %}
}
{% end %}