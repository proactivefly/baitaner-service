package user

import (
	"encoding/json"
	"fmt"
	"gofly/global"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/**
*使用 Index 是省略路径中的index
*本路径为： /admin/user/login -省去了index
 */
type Index struct{}

func init() {
	fPath := Index{}
	gf.Register(&fPath, reflect.TypeOf(fPath).PkgPath())
}

// GetUserInfoByOpenid 通过openid获取用户信息
func (api *Index) GetUserInfoByOpenid(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	// 详情见文档 https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	ref := gf.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code", global.App.Config.WXConf.AppID, global.App.Config.WXConf.AppSecret, code))
	var parameter map[string]interface{}
	if err := json.Unmarshal([]byte(ref), &parameter); err == nil {
		rootUrl, _ := model.DB().Table("common_config").Where("keyname", "rooturl").Value("keyvalue")
		user, _ := model.DB().Table("stall_user").Where("openid", parameter["openid"]).First()
		if user != nil {
			//token
			token := getToken(user)
			// 获取头像 并断言为string 类型 .(string) 为类型断言
			if !strings.Contains(user["avatar"].(string), "http") && rootUrl != nil {
				user["avatar"] = rootUrl.(string) + user["avatar"].(string)
			}
			results.SuccessLogin(c, "直接获取已有的用户数据", user, token, nil)
		} else { //不存在则添加一条
			parameter["create_time"] = time.Now().Unix()
			parameter["avatar"] = "resource/staticfile/avatar.png" // 默认头像
			delete(parameter, "session_key")
			addId, err := model.DB().Table("stall_user").Data(parameter).InsertGetId()
			if err != nil {
				results.Failed(c, "添加账号失败", err)
			} else {
				//fmt.Sprintf("%s%v", "hl_", addId) 将字符串 "hl_" 和 addId 的值拼接在一起。
				//model.DB().Table("stall_user").Data(map[string]interface{}{"name": fmt.Sprintf("%s%v", "hl_", addId)}).Where("id", addId).Update()
				user, _ := model.DB().Table("stall_user").Where("id", addId).First()
				//token
				token := getToken(user)
				if !strings.Contains(user["avatar"].(string), "http") && rootUrl != nil {
					user["avatar"] = rootUrl.(string) + user["avatar"].(string)
				}
				results.SuccessLogin(c, "添加并获取token！", user, token, nil)
			}
		}
	} else {
		results.Failed(c, "获取openid失败", err)
	}
}

// 获取Token
func getToken(user map[string]interface{}) interface{} {
	token := middleware.GenerateToken(&middleware.UserClaims{
		ID:             user["id"].(int64),
		Openid:         user["openid"].(string),
		StandardClaims: jwt.StandardClaims{},
	})
	return token
}

// GetUserinfo 获取用户信息
func (api *Index) GetUserinfo(c *gin.Context) {
	//当前用户
	token := c.Request.Header.Get("Authorization")
	user := middleware.ParseToken(token)
	data, err := model.DB().Table("stall_user").Where("id", user.ID).First()
	if err != nil {
		results.Failed(c, "获取用户信息失败", err)
	} else {
		results.Success(c, "获取用户信息成功！", data, nil)
	}
}

// UpdateUserInfo 保存
func (api *Index) UpdateUserInfo(c *gin.Context) {
	//获取post传过来的data
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	token := c.Request.Header.Get("Authorization")
	user := middleware.ParseToken(token)
	_ = json.Unmarshal(body, &parameter)
	res, err := model.DB().Table("stall_user").
		Data(parameter).
		Where("id", user.ID).
		Update()
	if err != nil {
		results.Failed(c, "更新失败", err)
	} else {
		results.Success(c, "更新成功！", res, nil)
	}
}
