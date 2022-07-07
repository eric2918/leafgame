package internal

import (
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo account inputs", commandEcho)
}

func commandEcho(args []interface{}) interface{} {
	fmt.Println(args)
	return fmt.Sprintf("%v", args)
}
