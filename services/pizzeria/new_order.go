package pizzeria

import app "tonky/holistic/apps/pizzeria"

type NewOrder struct {
	Content string
}

func (no NewOrder) ToApp() app.NewOrder {
	return app.NewOrder{Content: no.Content}
}
