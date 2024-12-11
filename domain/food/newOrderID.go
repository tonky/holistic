package food

import "github.com/google/uuid"

func NewOrderID(us string) (OrderID, error) {
	id, err := uuid.Parse(us)
	if err != nil {
		return OrderID{}, err
	}

	return OrderID{id: id}, nil
}

func (o OrderID) String() string {
	return o.id.String()
}
