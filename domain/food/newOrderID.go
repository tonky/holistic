package food

import "github.com/google/uuid"

func NewOrderID(us string) (OrderID, error) {
	id, err := uuid.Parse(us)
	if err != nil {
		return OrderID{}, err
	}

	return OrderID{ID: id}, nil
}

func (o OrderID) String() string {
	return o.ID.String()
}
