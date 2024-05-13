package stall

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
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
		MDB.Where("s.status", "like", "%"+stallStatus+"%")
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

// GetDetail 获取内容
func (api *Index) GetDetail(c *gin.Context) {
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

// 删除
func (api *Index) Delete(c *gin.Context) {

	/**
	这段 Go 代码是在从一个 HTTP 请求中读取其主体内容（Body）。让我逐步为你解释：

	c.Request.Body:

	c 可能是一个 HTTP 处理程序的上下文对象，例如，在 Gin、Echo 或标准库的 http.Handler 中。
	c.Request 是该上下文对象中的 HTTP 请求对象。
	c.Request.Body 是该 HTTP 请求的主体部分，它是一个 io.ReadCloser 接口。io.ReadCloser 是一个结合了 io.Reader 和 io.Closer 的接口，这意味着你可以从它读取数据，并在完成后关闭它。
	io.ReadAll(c.Request.Body):

	io.ReadAll 是一个 Go 标准库中的函数，它读取 io.Reader 接口的所有数据，并返回读取的字节切片和可能遇到的任何错误。
	在这里，我们传递 c.Request.Body 作为参数，意味着我们想读取整个 HTTP 请求的主体内容。
	body, _ := ...:

	这是一个 Go 中的简短变量声明和初始化。body 是一个字节切片（[]byte），它存储了从 c.Request.Body 读取的所有数据。
	_ 是 Go 中的一个特殊标识符，用于忽略返回值。在这里，我们忽略了可能由 io.ReadAll 返回的错误。通常，在实际的应用程序中，我们不建议忽略错误，而是应该检查并处理它，例如使用 log.Fatal 记录错误或返回给调用者。
	综上所述，这段代码的主要目的是从 HTTP 请求中读取主体内容，并将其存储在 body 字节切片中。但是，它忽略了可能发生的任何错误，这在实际应用中可能不是一个好的做法。
	*/
	body, _ := io.ReadAll(c.Request.Body)
	/*
		这段 Go 代码定义了一个名为 parameter 的变量，其类型为 map[string]interface{}。

		让我们逐步解释这段代码：

		var 关键字：这是 Go 语言中用于声明变量的关键字。

		parameter：这是变量的名称。

		map[string]interface{}：这是变量的类型。它定义了一个映射（或字典），其中键是 string 类型，而值是 interface{} 类型。

		string：表示映射中的键是字符串类型。
		interface{}：是 Go 语言中的一个空接口，它可以表示任何类型。因此，这意味着映射的值可以是任何类型。
		所以，parameter 是一个可以存储任何类型值的映射，其键是字符串。

	*/
	var parameter map[string]interface{}
	// 打印方法
	//fmt.Println("___", string(body))
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"]
	res2, err := model.DB().Table("stall").Where("id", id).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", res2, nil)
	}
}

// 删除
func (api *Index) UpdateStatus(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	// 打印方法
	//fmt.Println("___", string(body))
	_ = json.Unmarshal(body, &parameter)
	id := parameter["id"]
	status := parameter["status"]
	res2, err := model.DB().Table("stall").Data(map[string]interface{}{"status": status}).Where("id", id).Update()
	if err != nil {
		results.Failed(c, "更新失败", err)
	} else {
		results.Success(c, "更新成功！", res2, nil)
	}
}
