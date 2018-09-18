package messages

type StatMessage struct {
	Memory float64 `json:"memory"`
	Cpu    float64 `json:"cpu"`
}

func (m StatMessage) Key() string {
	return "stat"
}
