package services

import (
	"strings"
	"tonky/holistic/generator/domain"

	"github.com/open2b/scriggo/builtin"
)

type TopicDesc struct {
	Name         string
	TopicName    string
	DomainObject domain.Object
}

func (td TopicDesc) InterfaceName() string {
	return "I" + builtin.Capitalize(td.Name)
}

func (td TopicDesc) StructName() string {
	return builtin.Capitalize(td.Name)
}

func (td TopicDesc) ModelName() string {
	return td.DomainObject.Domain + "." + td.DomainObject.Name
}

func (td TopicDesc) SnakeFileName() string {
	return strings.Replace(td.TopicName, ".", "_", -1)
}

func (td TopicDesc) AppVarName() string {
	return td.StructName()
}

var TopicFoodOrderCreated = TopicDesc{Name: "foodOrderCreated", TopicName: "food.order.created", DomainObject: domain.FoodOrder}
var TopicFoodOrderUpdated = TopicDesc{Name: "foodOrderUpdated", TopicName: "food.order.updated", DomainObject: domain.FoodOrder}
var TopicFoodOrderPaid = TopicDesc{Name: "foodOrderPaid", TopicName: "food.order.paid", DomainObject: domain.AccountingOrder}
var TopicFoodOrderShipped = TopicDesc{Name: "foodOrderShipped", TopicName: "food.order.shipped", DomainObject: domain.AccountingOrder}

var kafkaTopics = []TopicDesc{TopicFoodOrderCreated, TopicFoodOrderUpdated, TopicFoodOrderPaid, TopicFoodOrderShipped}

func KafkaTopics() []TopicDesc {
	return kafkaTopics
}
