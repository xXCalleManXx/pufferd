package errors

type ServerOffline struct {
}

func (e ServerOffline) Error() string {
	return "Server offline"
}

func NewServerOffline() error {
	return ServerOffline{}
}
