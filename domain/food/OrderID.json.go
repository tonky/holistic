package food

import (
	"encoding/json"
	"fmt"
)

func (o *OrderID) UnmarshalJSON(b []byte) error {
	fmt.Println("OrderID.UnmarshalJSON: ", string(b[:]))

	oid, err := NewOrderID(string(b))
	if err != nil {
		return err
	}

	*o = oid

	return nil
}

func (o OrderID) MarshalJSON() ([]byte, error) {
	fmt.Println("OrderID.MarshalJSON: ", o.String())

	return json.Marshal(o.String())
}
