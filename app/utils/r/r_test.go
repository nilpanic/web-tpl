package r_test

import (
	"fmt"
	"testing"
	"time"

	"web-tpl/app/utils/r"
)

func TestGo(t *testing.T) {
	for i := 0; i < 5; i++ {
		r.Go(func(val any) {

		}, i)

	}
	time.Sleep(time.Second * 1)
	fmt.Println("hello")
}
