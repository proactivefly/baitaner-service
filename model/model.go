package model

import (
	"database/sql"
	"fmt"
	"gofly/global"
	"time"

	"gofly/utils/gform" //数据库操作

	_ "github.com/go-sql-driver/mysql"
)

var err error
var engin *gform.Engin

// MyInit 取得数据库连接实例 参数 starType 为接口类型，表示连接数据库的类型
func MyInit(starType interface{}) {
	// 输出日志，提示正在连接数据库
	global.App.Log.Info(fmt.Sprintf("连接数据库中:%v", starType))
	// 初始化配置
	global.App.Config.InitializeConfig()
	// 构建数据库连接字符串
	dsbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=1000ms", global.App.Config.DBconf.Username, global.App.Config.DBconf.Password, global.App.Config.DBconf.Hostname, global.App.Config.DBconf.Hostport, global.App.Config.DBconf.Database)
	// 打开数据库连接
	engin, err = gform.Open(&gform.Config{Driver: global.App.Config.DBconf.Driver, Dsn: dsbSource, Prefix: global.App.Config.DBconf.Prefix})
	// 判断连接是否成功
	if err != nil {
		// 输出日志，提示数据库连接实例错误
		global.App.Log.Info(fmt.Sprintf("数据库连接实例错误: %v", err))
	} else {
		// 输出日志，提示连接数据库成功
		global.App.Log.Info(fmt.Sprintf("连接数据库成功:%v", starType))
		// 设置连接池最大空闲连接数为10
		engin.GetExecuteDB().SetMaxIdleConns(10) //连接池最大空闲连接数,不设置, 默认无
		// 设置连接池最大连接数为50
		engin.GetExecuteDB().SetMaxOpenConns(50) // 连接池最大连接数,不设置, 默认无限
		// 设置连接池连接的最大生命周期为59秒
		engin.GetExecuteDB().SetConnMaxLifetime(59 * time.Second) //时间比超时时间短
		// 执行SQL语句，设置SQL模式为'NO_ENGINE_SUBSTITUTION'
		engin.GetQueryDB().Exec("SET @@sql_mode='NO_ENGINE_SUBSTITUTION';")
	}
}

// DB controller层调用
func DB() gform.IOrm {
	return engin.NewOrm()
}
func DBEV() *gform.Engin {
	return engin
}

// CreateDataBase 新建数据库
func CreateDataBase(Username, Password, Hostname, HostPort, Database interface{}) {
	global.App.Config.InitializeConfig()
	dsbSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=1000ms", Username, Password, Hostname, HostPort, "")
	engin, err = gform.Open(&gform.Config{Driver: global.App.Config.DBconf.Driver, Dsn: dsbSource})
	if err != nil {
		global.App.Log.Info(fmt.Sprintf("创建时，数据库连接实例错误: %v", err))
	} else {
		engin.GetQueryDB().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci", Database))
	}
}

// ExecSql 导入数据库文件
func ExecSql(rows string) (sql.Result, error) {
	Result, err := engin.GetExecuteDB().Exec(rows)
	if err != nil {
		global.App.Log.Info(fmt.Sprintf("导入数据失败:%v。%v", err, Result))
		return nil, err
	}
	return Result, nil
}

// GetTotal 取得总行数
func GetTotal(tableName string, wheres map[string]interface{}) int64 {
	total, _ := DB().Table(tableName).Where(wheres).Count()
	return total
}
