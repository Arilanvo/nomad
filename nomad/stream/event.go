package stream

const (
	AllKeys = "*"
)

type Topic string

type Event struct {
	Topic      Topic
	Type       string
	Key        string
	FilterKeys []string
	Index      uint64
	Payload    interface{}
}
