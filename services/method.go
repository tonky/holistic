package services

type Method string

const (
	Create Method = "create"
	Read   Method = "read"
	Update Method = "update"
	Delete Method = "delete"
)

func (m Method) HttpName() string {
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
