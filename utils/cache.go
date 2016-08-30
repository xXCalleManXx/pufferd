package utils

type cache struct {
	buffer []string
	capacity int
}

type Cache interface {
	Read() []string

	Write (b []byte) (n int, err error)
}

func CreateCache() cache {
	capacity := 50
	return cache{
		buffer: make([]string, 0),
		capacity: capacity,
	}
}

func (c cache) Read() []string {
	return c.buffer
}

func (c cache) Write(b []byte) (n int, err error) {
	msg := string(b)
	if len(c.buffer) == c.capacity {
		newbuffer := make([]string, c.capacity - 1)
		copy(newbuffer, c.buffer)
		c.buffer = newbuffer
	}
	c.buffer = append(c.buffer, msg)
	n = len(b)
	return
}