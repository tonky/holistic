package describer

type MethodAction string

const (
	Create MethodAction = "create"
	Read   MethodAction = "read"
	Update MethodAction = "update"
	Delete MethodAction = "delete"
)

func (m MethodAction) HttpName() string {
	switch m {
	case Create:
		return "Post"
	case Read:
		return "Get"
	case Update:
		return "Put"
	case Delete:
		return "Delete"
	}

	return "!!!!!!! undefined"
}
