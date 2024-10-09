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
		return "POST"
	case Read:
		return "GET"
	case Update:
		return "PUT"
	case Delete:
		return "DELETE"
	}

	return "!!!!!!! undefined"
}
