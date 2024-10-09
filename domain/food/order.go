package food

import (
	"hf/holistic/domain"
)

var OrderID = domain.Object{
	Fields: map[domain.FieldName]domain.Kind{
		"id": domain.UUID,
	},
}

var Order = domain.Object{
	Fields: map[domain.FieldName]domain.Kind{
		"id":      domain.UUID,
		"content": domain.String,
	},
}
