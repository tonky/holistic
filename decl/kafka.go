package decl

import (
	"tonky/holistic/describer"
)

var TopicFoodOrderCreated = describer.TopicDesc{Name: "foodOrderCreated", TopicName: "food.order.created", DomainObject: FoodOrder}
var TopicFoodOrderUpdated = describer.TopicDesc{Name: "foodOrderUpdated", TopicName: "food.order.updated", DomainObject: FoodOrder}
var TopicAccountingOrderPaid = describer.TopicDesc{Name: "accountingOrderPaid", TopicName: "accounting.order.paid", DomainObject: AccountingOrder}
var TopicShippingOrderShipped = describer.TopicDesc{Name: "shippingOrderShipped", TopicName: "shipping.order.shipped", DomainObject: ShippedOrder}

var KafkaTopics = []describer.TopicDesc{
	TopicFoodOrderCreated,
	TopicFoodOrderUpdated,
	TopicAccountingOrderPaid,
	TopicShippingOrderShipped,
}
