package pizzeria

import (
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/domain/food"
)

type NewOrder struct {
	Content string
}

type UpdateOrder struct {
	ID      food.OrderID
	Content string
	IsFinal bool
}

func (no NewOrder) ToApp() app.NewOrder {
	return app.NewOrder{Content: no.Content}
}

func (uo UpdateOrder) ToApp() app.UpdateOrder {
	return app.UpdateOrder{OrderID: uo.ID, Content: uo.Content, IsFinal: uo.IsFinal}
}
