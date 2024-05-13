package controller

import (
	"fmt"
	_ "gofly/app/baitaner/product"
	_ "gofly/app/baitaner/stall"
	_ "gofly/app/baitaner/user"
)

func init() {
	fmt.Println("init controller")
}
