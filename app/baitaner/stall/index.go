package stall

import (
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

/**
*使用 Index 是省略路径中的index
*本路径为： /admin/user/login -省去了index
 */
func init() {
	gf.Register(&Index{}, reflect.TypeOf(Index{}).PkgPath())
}

type Index struct {
}

// 获取商店列表
func (api *Index) Get_list(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := model.DB().Table("stall").Fields("id,name,logo,status,owner_uid,star,create_at")
	if name != "" {
		MDB.Where("name", "like", "%"+name+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("created_at asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		//处理结果集
		//for _, val := range list {
		//	catename, _ := model.DB().Table("business_store_cate").Where("id", val["cid"]).Value("name")
		//	val["catename"] = catename
		//}
		var totalCount int64
		totalCount, _ = model.DB().Table("stall").Count()
		results.Success(c, "摊位列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

// 获取内容
func (api *Index) Get_detail(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		results.Failed(c, "请传参数id", nil)
	} else {
		res2, err := model.DB().Table("stall").Where("id", id).First()
		if err != nil {
			results.Failed(c, "获取内容失败", err)
		} else {
			results.Success(c, "获取内容成功！", res2, nil)
		}
	}
}
