package stall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"reflect"
	"strconv"
	"time"
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

// GetList 获取商店列表
func (api *Index) GetList(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	stallTypeId := c.DefaultQuery("type_id", "")
	stallStatus := c.DefaultQuery("status", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//通过左连接查询
	MDB := model.DB().Table("stall s").LeftJoin("common_dictionary_stall_type cds on s.type_id = cds.id")
	MDB.Fields("s.id", "s.name", "s.type_id", "s.logo", "s.owner_uid", "s.created_at", "s.star", "s.status", "s.description", "s.owner_phone", "s.wechat_id", "s.owner_uid", "s.logo", "cds.keyname type_name")
	// 查询stall表中type_id对应的
	if name != "" {
		MDB.Where("name", "like", "%"+name+"%")
	}
	if stallStatus != "" {
		MDB.Where("status", "like", "%"+name+"%")
	}
	if stallTypeId != "" {
		MDB.Where("type_id", "like", "%"+stallTypeId+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("created_at desc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		//处理结果集
		for _, val := range list {
			//取用户名
			owenName, _ := model.DB().Table("stall_users").Where("id", val["owner_uid"]).Value("username")
			val["owner_name"] = owenName

			//把时间转换为
			//把val["created_at"] 转成int(64)
			//fmt.Printf("什么类型呢%T", val["created_at"])
			//time64 := val["created_at"].Unix()

			// 检查"created_at"的类型
			timeVal, ok := val["created_at"].(time.Time)
			if ok {
				time64 := timeVal.Unix()
				val["created_at"] = gf.TimestampString(time64, "datetime")
			}

			//if createdAt, ok := val["created_at"].(time.Time); ok {
			//	time64 := createdAt.Unix()
			//	val["created_at"] = gf.TimestampString(time64, "datetime")
			//}
		}
		var totalCount int64
		totalCount, _ = model.DB().Table("stall").Count()
		results.Success(c, "摊位列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list}, nil)
	}
}

// Get_detail 获取内容
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

// Delete 删除内容
func (api *Index) Delete(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	fmt.Println("——————————", c)
	if id == "" {
		results.Failed(c, "请传参数id", nil)
	} else {
		res, err := model.DB().Table("stall").Where("id", id).Delete()
		if err != nil {
			results.Failed(c, "删除失败", err)
		} else {
			results.Success(c, "删除成功！", res, nil)
		}
	}
}
