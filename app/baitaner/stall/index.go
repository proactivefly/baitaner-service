package stall

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

// GetStallRecommendListRecommend 获取首页推荐类表
func (api *Index) GetStallRecommendListRecommend(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	_pageSize := c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//通过左连接查询
	MDB := model.DB().Table("stall").Fields("id", "name", "logo", "description", "type_id", "address", "lon", "lat")
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("created_at asc").Get()
	if err != nil {
		results.Failed(c, err.Error(), nil)
	} else {
		//解析经纬度 获取距离
		for _, v := range list {
			if v["lon"] != nil {
				//v["distance"] = gf.GetDistance(gf.Location{Lon: 116.403875, Lat: 39.92178}, gf.Location{Lon: v["lon"].(float64), Lat: v["lat"].(float64)})
				v["distance"] = "80.8m"
			} else {
				v["distance"] = ""
			}
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

// GetStallType 获取类型接口
func (api *Index) GetStallType(c *gin.Context) {
	res, err := model.DB().Table("common_dictionary_stall_type").Fields("keyname text", "keyvalue value").Order("value asc").Get()
	if err != nil {
		results.Failed(c, "获取类型失败", err)
	} else {
		results.Success(c, "获取类型！", res, nil)
	}
}

// Save 保存
func (api *Index) Save(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	//json.Unmarshal: 这是 json 包中提供的一个函数，用于将 JSON 编码的数据解码（或反序列化）为 Go 语言中的数据结构。
	_ = json.Unmarshal(body, &parameter)
	var ID float64 = 0
	if parameter["id"] != nil {
		ID = parameter["id"].(float64)
	}
	stallName := parameter["name"]
	if stallName == nil {
		results.Failed(c, "请填写摊位名称", nil)
	}
	stallType := parameter["type_id"]
	if stallType == nil {
		results.Failed(c, "请选择摊位类型", nil)
	}
	stallAddress := parameter["address"]
	if stallAddress == nil {
		results.Failed(c, "请填写摊位地址", nil)
	}
	stallLon := parameter["lon"]
	if stallLon == nil {
		results.Failed(c, "请填写摊位经度", nil)
	}
	stallLat := parameter["lat"]
	if stallLat == nil {
		results.Failed(c, "请填写摊位纬度", nil)
	}
	stallDescription := parameter["description"]
	if stallDescription == nil {
		results.Failed(c, "请填写摊位描述", nil)
	}
	wechat_id := parameter["wechat_id"]
	if wechat_id == nil {
		results.Failed(c, "请填写微信号", nil)
	}
	logo := parameter["logo"]
	if logo == nil {
		results.Failed(c, "请填写摊位logo", nil)
	}
	owner_phone := parameter["owner_phone"]
	//正则验证手机号
	if owner_phone == nil {
		results.Failed(c, "请填写手机号", nil)
	} else if !regexp.MustCompile(`^1[3456789]\d{9}$`).MatchString(owner_phone.(string)) {
		results.Failed(c, "请填写正确手机号", nil)
	}
	// 验证手机号是否存在
	res, err := model.DB().Table("stall").Where("mobile", owner_phone.(string)).First()
	if err != nil {
		results.Failed(c, "获取内容失败", err)
	} else if res["id"] != ID {
		results.Failed(c, "手机号已存在", nil)
	}

	if ID == 0 { // 新增
		delete(parameter, "id")
		parameter["created_at"] = time.Now().Unix()
		parameter["status"] = "closed"
		addId, err := model.DB().Table("stall").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加失败", err)
		} else {
			results.Success(c, "添加成功！", addId, nil)
		}
	} else { // 更新
		res, err := model.DB().Table("stall").
			Data(parameter).
			Where("id", ID).
			Update()
		if err != nil {
			results.Failed(c, "更新失败", err)
		} else {
			results.Success(c, "更新成功！", res, nil)
		}
	}
}
