package operations

import "fmt"

type Download struct {
	File string
}

func (d *Download) Run() {
	fmt.Println("Downloading file: " + d.File)
}
