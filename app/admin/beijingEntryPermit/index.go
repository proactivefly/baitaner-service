package beijingEntryPermit

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strconv"
	"time"
)

func init() {
	path := Index{}
	gf.Register(&path, reflect.TypeOf(path).PkgPath())
}

// Index 用于自动注册路由
type Index struct {
}

func (api *Index) GetList(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	serviceCycle := c.DefaultQuery("service_cycle", "")
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	status := c.DefaultQuery("status", "")
	pageNo, _ := strconv.Atoi(page) // 字符串转int
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := model.DB().Table("beijing_entry_permit").Fields("*")
	if name != "" {
		MDB.Where("name", "like", "%"+name+"%")
	}
	if serviceCycle != "" {
		MDB.Where("service_cycle", serviceCycle)
	}
	if status != "" {
		MDB.Where("status", status)
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		for i := range list {
			item := list[i]

			// 处理 start_time
			if val, ok := item["start_time"].(time.Time); ok {
				item["start_time"] = val.Format(time.DateTime)
			}

			// 处理 end_time
			if val, ok := item["end_time"].(time.Time); ok {
				item["end_time"] = val.Format(time.DateTime)
			}

			// 处理 serviceCycle
			var serviceCycleText string
			fmt.Printf("值: %v, reflect类型: %v, T类型: %T\n", item["service_cycle"], reflect.TypeOf(item["service_cycle"]), item["service_cycle"])
			if v, ok := item["service_cycle"].(int64); ok {
				switch v {
				case 1:
					serviceCycleText = "月"
				case 2:
					serviceCycleText = "季度"
				case 3:
					serviceCycleText = "年"

				default:
					serviceCycleText = "未知"
				}
			}

			item["service_cycle_text"] = serviceCycleText
		}
		var totalCount int64
		totalCount, _ = model.DB().Table("beijing_entry_permit").Count()
		results.Success(c, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// Save 保存编辑
func (api *Index) Save(c *gin.Context) {
	// 获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	// 校验必填字段
	requiredFields := []string{"auth_key", "fangtang_key", "service_cycle", "start_time"}
	if !gf.ValidateRequiredFields(c, parameter, requiredFields) {
		return
	}
	// 当前用户
	//getUser, _ := c.Get("user")
	//user := getUser.(*middleware.UserClaims)
	//parameter["create_by"] = user.Name
	var itemId float64 = 0
	if parameter["id"] != nil {
		itemId = parameter["id"].(float64)
	}
	if parameter["service_cycle"] != nil {
		parameter["service_cycle"] = parameter["service_cycle"]
	}
	if parameter["start_time"] != nil {
		parameter["start_time"] = parameter["start_time"]
	}
	// 计算 end_time
	var endTime time.Time
	startTimeFloat, ok1 := parameter["start_time"].(float64)       // 如果是 timestamp 数字（常见于前端传值）
	serviceCycleFloat, ok2 := parameter["service_cycle"].(float64) // 同样处理 float64
	if ok1 && ok2 {
		startTime := time.Unix(int64(startTimeFloat), 0)
		switch int(serviceCycleFloat) {
		case 1:
			endTime = startTime.AddDate(0, 1, 0) // 加一个月
		case 2:
			endTime = startTime.AddDate(0, 3, 0) // 加一个季度
		case 3:
			endTime = startTime.AddDate(1, 0, 0) // 加一年
		default:
			endTime = startTime // 默认不变
		}
		parameter["end_time"] = endTime.Unix()
	}
	//
	parameter["created_at"] = time.Now().Unix()
	if itemId == 0 { // 新增逻辑
		delete(parameter, "id")
		addId, err := model.DB().Table("beijing_entry_permit").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加失败", err)
		} else {
			results.Success(c, "添加成功！", addId, nil)
		}
	} else { // 修改逻辑
		res, err := model.DB().Table("beijing_entry_permit").
			Data(parameter).
			Where("id", itemId).
			Update()
		if err != nil {
			results.Failed(c, "更新失败", err)
		} else {
			results.Success(c, "更新成功！", res, nil)
		}
	}
}
