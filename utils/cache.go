package utils

type cache struct {
	buffer []string
	capacity int
}

type Cache interface {
	Read() []string

	Write (b []byte) (n int, err error)
}

func CreateCache() *cache {
	capacity := 3
	return &cache{
		buffer: make([]string, 0),
		capacity: capacity,
	}
}

func (c *cache) Read() []string {
	result := make([]string, len(c.buffer))
	for k, v := range c.buffer {
		result[k] = v
	}
	return result
}

func (c *cache) Write(b []byte) (n int, err error) {
	if len(c.buffer) == c.capacity {
		newbuffer := make([]string, c.capacity - 1)
		copy(newbuffer, c.buffer)
		c.buffer = newbuffer
	}
	c.buffer = append(c.buffer, string(b))
	n = len(b)
	return
}