package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"trojan-panel-core/api"
	"trojan-panel-core/app"
	"trojan-panel-core/core"
	"trojan-panel-core/dao"
	"trojan-panel-core/dao/redis"
	"trojan-panel-core/middleware"
	"trojan-panel-core/model/constant"
	"trojan-panel-core/router"
)

func main() {
	serverConfig := core.Config.ServerConfig
	r := gin.Default()
	router.Router(r)
	_ = r.Run(fmt.Sprintf(":%d", serverConfig.Port))
	defer closeResource()
}

func init() {
	core.InitConfig()
	middleware.InitLog()
	
	// 初始化数据库连接
	initDatabase()
	
	// 初始化中间件
	middleware.InitCron()
	middleware.InitRateLimiter()
	middleware.InitTrafficControl()
	
	// 初始化API和应用
	api.InitValidator()
	api.InitGrpcServer()
	app.InitApp()
}

// 初始化数据库连接
func initDatabase() {
	// 初始化MySQL
	if err := dao.InitMySQL(); err != nil {
		logrus.Errorf("初始化MySQL失败: %v", err)
	}
	
	// 初始化SQLite
	if err := dao.InitSqlLite(); err != nil {
		logrus.Errorf("初始化SQLite失败: %v", err)
	}
	
	// 初始化Redis
	if err := redis.InitRedis(); err != nil {
		logrus.Errorf("初始化Redis失败: %v", err)
	}
}
func closeResource() {
	dao.CloseDb()
	dao.CloseSqliteDb()
	redis.CloseRedis()
}
