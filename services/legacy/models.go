package legacy

import (
	"fmt"
	app "tonky/holistic/apps/legacy"
	"tonky/holistic/domain/food"
)

type OrderID struct {
	ID string
}

func (o OrderID) ToApp() (food.OrderID, error) {
	return food.NewOrderID(o.ID)
}

type NewOrder struct {
	Content string
}

func (no NewOrder) ToApp() (app.NewOrder, error) {
	return app.NewOrder{Content: no.Content}, nil
}

type UpdateOrder struct {
	ID      string
	Content string
	IsFinal bool
}

func (uo UpdateOrder) ToApp() (app.UpdateOrder, error) {
	oid, err := food.NewOrderID(uo.ID)
	if err != nil {
		return app.UpdateOrder{}, err
	}

	return app.UpdateOrder{ID: oid, Content: uo.Content, IsFinal: uo.IsFinal}, nil
}

var ErrEmptyContent = fmt.Errorf("empty content")

func (no NewOrder) Validate() error {
	if no.Content == "" {
		return ErrEmptyContent
	}

	return nil
}
