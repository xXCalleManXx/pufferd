package messages

type ConsoleMessage struct {
	Logs []string `json:"logs"`
}

func (m ConsoleMessage) Key() string {
	return "console"
}
