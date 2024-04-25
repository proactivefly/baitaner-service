package controller

import (
	"fmt"
	_ "gofly/app/baitaner/product"
	_ "gofly/app/baitaner/stall"
)

func init() {
	fmt.Println("init controller")
}
