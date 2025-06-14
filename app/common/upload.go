package common

import (
	"fmt"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	fPath := Upload{}
	gf.Register(&fPath, reflect.TypeOf(fPath).PkgPath())
}

type Upload struct {
}

// Image 1.上传单文件
func (api *Upload) Image(c *gin.Context) {
	fmt.Println("上传图片接口")
	cid := c.DefaultPostForm("cid", "1")
	// 单个文件
	file, err := c.FormFile("file")
	if err != nil {
		results.Failed(c, "获取数据失败，", err)
		return
	}
	nowTime := time.Now().Unix() //当前时间
	getUser, _ := c.Get("user")  //取值 实现了跨中间件取值
	user := getUser.(*middleware.UserClaims)
	//时间查询-获取当天时间
	dayTime := time.Now().Format("20060102")
	//文件唯一性，拼接文件名
	fileUniName := fmt.Sprintf("%s%s%v", file.Filename, dayTime, user.ID)
	sha1Str := gf.Md5(fileUniName)
	//开始
	dayStar, _ := time.Parse("20060102", dayTime+" 00:00:00")
	dayStarTimes := dayStar.Unix() //时间戳
	//结束
	dayEnd, _ := time.Parse("20060102", dayTime+" 23:59:59")
	dayEndTimes := dayEnd.Unix() //时间戳
	rootUrl, _ := model.DB().Table("common_config").Where("keyname", "rootUrl").Value("keyvalue")
	// 判断文件是否已经存在
	attachment, _ := model.DB().Table("attachment").Where("uid", user.ID).
		WhereBetween("uploadtime", []interface{}{dayStarTimes, dayEndTimes}).
		Where("sha1", sha1Str).Fields("id,title,url").First()
	if attachment != nil { //文件是否已经存在
		c.JSON(200, gin.H{
			"id":       attachment["id"],
			"uid":      sha1Str,
			"name":     attachment["name"],
			"status":   "done",
			"url":      fmt.Sprintf("%s%s", rootUrl, attachment["url"]),
			"response": "文件已上传",
			"time":     nowTime,
		})
		c.Abort()
		return
	}

	//fmt.Sprintf函数根据格式字符串和参数生成一个新的字符串，但不会打印出来。格式字符串"%s%s%s"定义了三个字符串插值的位置，它们将由后面的参数填充。
	filePath := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(filePath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filePath, os.ModePerm)
		}
	}
	//上传到的路径
	filenameArr := strings.Split(file.Filename, ".")
	nameStr := gf.Md5(fmt.Sprintf("%v%s", nowTime, filenameArr[0])) //组装文件保存名字
	//path := 'resource/uploads/20060102150405test.xlsx'
	fileFilename := fmt.Sprintf("%s%s%s", nameStr, ".", filenameArr[1]) //文件加.后缀
	path := filePath + fileFilename
	// 上传文件到指定的目录
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(200, gin.H{
			"uid":      sha1Str,
			"name":     file.Filename,
			"status":   "error",
			"response": "上传失败",
			"time":     nowTime,
		})
	} else {
		//保存数据
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		insertData := map[string]interface{}{
			"accountID":  user.Accountid,
			"cid":        cid,
			"uid":        user.ID,
			"sha1":       sha1Str,
			"title":      filenameArr[0],
			"name":       file.Filename,
			"url":        path,
			"storage":    dir + strings.Replace(path, "/", "\\", -1),
			"uploadtime": nowTime,
			"updatetime": nowTime,
			"filesize":   file.Size,
			"mimetype":   file.Header["Content-Type"][0],
		}
		fileId, _ := model.DB().Table("attachment").Data(insertData).InsertGetId()
		c.JSON(200, gin.H{
			"id":       fileId,
			"uid":      sha1Str,
			"name":     file.Filename,
			"status":   "done",
			"url":      fmt.Sprintf("%s%s", rootUrl, path),
			"thumb":    path,
			"response": "上传成功",
			"time":     nowTime,
		})
	}
}

// File 1.文件
func (api *Upload) File(c *gin.Context) {
	cid := c.DefaultPostForm("cid", "1")
	// 单个文件
	file, err := c.FormFile("file")
	if err != nil {
		results.Failed(c, "获取数据失败，", err)
		return
	}
	nowTime := time.Now().Unix() //当前时间
	getuser, _ := c.Get("user")  //取值 实现了跨中间件取值
	user := getuser.(*middleware.UserClaims)
	//时间查询-获取当天时间
	day_time := time.Now().Format("2006-01-02")
	//文件唯一性
	file_uniname := fmt.Sprintf("%s%s%v", file.Filename, day_time, user.ID)
	sha1_str := gf.Md5(file_uniname)
	//开始
	day_star, _ := time.Parse("2006-01-02 15:04:05", day_time+" 00:00:00")
	day_star_times := day_star.Unix() //时间戳
	//结束
	day_end, _ := time.Parse("2006-01-02 15:04:05", day_time+" 23:59:59")
	day_end_times := day_end.Unix() //时间戳
	attachment, _ := model.DB().Table("attachment").Where("uid", user.ID).
		WhereBetween("uploadtime", []interface{}{day_star_times, day_end_times}).
		Where("sha1", sha1_str).Fields("id,title,url").First()
	if attachment != nil { //文件是否已经存在
		c.JSON(200, gin.H{
			"id":       attachment["id"],
			"uid":      sha1_str,
			"name":     attachment["name"],
			"status":   "done",
			"url":      attachment["url"],
			"response": "文件已上传",
			"time":     nowTime,
		})
		c.Abort()
		return
	}
	file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(file_path, os.ModePerm)
		}
	}
	//上传到的路径
	filename_arr := strings.Split(file.Filename, ".")
	name_str := gf.Md5(fmt.Sprintf("%v%s", nowTime, filename_arr[0])) //组装文件保存名字
	//path := 'resource/uploads/20060102150405test.xlsx'
	file_Filename := fmt.Sprintf("%s%s%s", name_str, ".", filename_arr[1]) //文件加.后缀
	path := file_path + file_Filename
	// 上传文件到指定的目录
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(200, gin.H{
			"uid":      sha1_str,
			"name":     file.Filename,
			"status":   "error",
			"response": "上传失败",
			"time":     nowTime,
		})
	} else {
		//保存数据
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		Insertdata := map[string]interface{}{
			"accountID":  user.Accountid,
			"cid":        cid,
			"uid":        user.ID,
			"sha1":       sha1_str,
			"title":      filename_arr[0],
			"name":       file.Filename,
			"url":        path,
			"storage":    dir + strings.Replace(path, "/", "\\", -1),
			"uploadtime": nowTime,
			"updatetime": nowTime,
			"filesize":   file.Size,
			"mimetype":   file.Header["Content-Type"][0],
		}
		model.DB().Table("attachment").Data(Insertdata).InsertGetId()
		rooturl, _ := model.DB().Table("common_config").Where("keyname", "rooturl").Value("keyvalue")
		c.JSON(200, gin.H{
			"code": 200,
			"data": map[string]interface{}{
				"download": fmt.Sprintf("%s%s", rooturl, path),
				"preview":  fmt.Sprintf("%s%s", rooturl, path),
				"url":      fmt.Sprintf("%s%s", rooturl, path),
			},
			"message": "上传成功",
		})
	}
}

// ThirdImage 编辑器保存第三方图片到本地
func (api *Upload) ThirdImage(c *gin.Context) {
	params, _ := gf.RequestParam(c)
	if url, ok := params["url"]; !ok || url == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"result": false,
			"data": map[string]interface{}{
				"url": "",
			},
			"message": "地址无效",
		})
	} else {
		file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
		if _, err := os.Stat(file_path); err != nil {
			if !os.IsExist(err) {
				os.MkdirAll(file_path, os.ModePerm)
			}
		}
		nowTime := time.Now().Unix() //当前时间
		localPicName := fmt.Sprintf("%vthir_%v", file_path, nowTime)
		imgtype, err := gf.DownPic(gf.InterfaceTostring(params["url"]), localPicName)
		rooturl, _ := model.DB().Table("common_config").Where("keyname", "rooturl").Value("keyvalue")
		c.JSON(200, gin.H{
			"code":    200,
			"result":  true,
			"err":     err,
			"status":  "done",
			"url":     fmt.Sprintf("%s%s%s", rooturl, localPicName, imgtype),
			"message": "上传成功",
		})
	}
}
