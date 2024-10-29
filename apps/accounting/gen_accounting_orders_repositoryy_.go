package accounting

/*
// compile-time check to make sure interface is implemented
var _ IAccountingOrdersRepository = new(OrdersRepository)

type IAccountingOrdersRepository interface {
	ReadOrderByFoodID(context.Context, food.OrderID) (accounting.Order, error)
}

type OrdersRepository struct {
	logger *logger.Slog
}

func NewAccountingOrdersRepository(logger *logger.Slog) *OrdersRepository {
	return &OrdersRepository{
		logger: logger,
	}
}
*/
