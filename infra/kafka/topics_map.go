package kafka

type Topic string

const (
	TopicFoodOrderCreated Topic = "food.order.created"
	TopicFoodOrderUpdated Topic = "food.order.updated"
)

var Topics struct {
	FoodOrderCreated Topic
}
