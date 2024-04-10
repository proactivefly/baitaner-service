package shop

import (
	"gofly/utils"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)
func init() {
	utils.Register(&Shop{}, reflect.TypeOf(Shop{}).PkgPath())
}
type Shop struct {
}

type Install struct {
}
// 安装页面
func (api *Install) Index(context *gin.Context){
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": "",
	})
}
// 获取商铺列表
func (api *Shop) Get_shoplist(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": []string{"1", "2", "3"},
	})
}