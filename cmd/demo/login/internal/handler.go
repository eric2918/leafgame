package internal

import (
	"reflect"
)

func handle(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

}
