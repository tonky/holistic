package food

import "github.com/google/uuid"

func NewOrderID(id string) (OrderID, error) {
	ui, err := uuid.Parse(id)

	if err != nil {
		return OrderID{}, err
	}

	return OrderID{ID: ui}, nil
}

func RandomOrderID() OrderID {
	return OrderID{ID: uuid.New()}
}

func (o OrderID) String() string {
	return o.ID.String()
}
