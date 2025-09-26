package testutil

type TestMessage struct {
	Content string `json:"content"`
}

func (tm TestMessage) Type() string {
	return "test_message"
}
