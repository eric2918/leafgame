package internal

import (
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo account inputs", echo)
}

func echo(args []interface{}) (ret interface{}, err error) {
	return fmt.Sprintf("%v", args), nil
}
