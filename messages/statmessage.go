package messages

type StatMessage struct {
	Memory int `json:"memory"`
	Cpu    int `json:"cpu"`
}

func (m StatMessage) Key() string {
	return "stat"
}
