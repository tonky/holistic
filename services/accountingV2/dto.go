package accountingV2

import (
	app "tonky/holistic/apps/accountingV2"
)

func serviceToAppNewFoodOrder(in NewFoodOrder) app.NewFoodOrder {
	return app.NewFoodOrder{
		Name:       in.Name,
		IsComplete: in.IsComplete,
	}
}
