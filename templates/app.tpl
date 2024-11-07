// AUTOGENERATED! DO NOT EDIT.

package {{ service.Name }}

import (
	{% if service.KafkaProducers %}
	"tonky/holistic/infra/kafkaProducer"
	{% end %}
	{% if service.KafkaConsumers %}
	"tonky/holistic/infra/kafkaConsumer"
	"context"
	{% end %}
	{% for d in client_deps %}
	"tonky/holistic/clients/{{ d.AppVarName() }}"
	{% end %}
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type App struct {
	Deps       do.Injector
	Logger     *logger.Slog

{% for ad in app_deps %}
    {{ cap(ad.AppVarName()) }} {{ ad.InterfaceName() }}
{% end %}
{% for d in client_deps.Dedup() %}
    {{ cap(d.AppVarName()) }} {{ d.AppVarName() }}.{{ d.InterfaceName() }}
{% end %}
}

func NewApp(deps do.Injector) (*App, error) {
	app := App{
		Deps:       deps,
		Logger:     do.MustInvoke[*logger.Slog](deps),
{% for ad in app_deps %}
        {{ cap(ad.AppVarName()) }}: do.MustInvokeAs[{{ ad.InterfaceName() }}](deps),
{% end %}
{% for d in client_deps.Dedup() %}
        {{ cap(d.AppVarName()) }}: do.MustInvokeAs[{{ d.AppVarName() }}.{{ d.InterfaceName() }}](deps),
{% end %}

	}

	return &app, nil
}

{% if service.KafkaConsumers %}
func (a App) RunConsumers() {
	a.Logger.Info(">> {{ service.Name}}.App.RunConsumers()")

	ctx := context.Background()
	{% for consumer in service.KafkaConsumers %}

	go func() {
		for err := range a.{{ cap(consumer.Name) }}Consumer.Run(ctx, a.{{ cap(consumer.Name) }}Processor) {
			a.Logger.Warn(err.Error())
		}
	}()
	{% end %}
}
{% end %}