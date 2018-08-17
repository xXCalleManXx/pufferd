package messages

type ConsoleMessage struct {
	Line string `json:"line"`
}

func (m ConsoleMessage) Key() string {
	return "console"
}
