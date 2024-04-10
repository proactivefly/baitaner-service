package controller

import (
	"fmt"
	_ "gofly/app/baitaner/product"
	_ "gofly/app/baitaner/shop"
)

func init(){
	fmt.Println("init controller")
}