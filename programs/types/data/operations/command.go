package operations

import "fmt"

type Command struct {
	Command string
}

func (c *Command) Run() {
	fmt.Println("Running command: " + c.Command)
}
