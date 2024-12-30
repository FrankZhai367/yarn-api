package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/config"
	"github.com/yarn/backend/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Options Request!"})
			return
		}

		//处理请求
		context.Next()

		if len(context.Errors) > 0 {
			context.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "server fail",
				"data": map[string]interface{}{},
			})
			context.Abort()
			return
		}
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	_, conf := config.LoadConf("config.yaml")
	db, _ := gorm.Open(mysql.Open(conf.Mysql.DataSourceName()), &gorm.Config{})

	r.Use(CrosHandler())
	service.NewUserService(db).RegisterHandler(r)
	service.NewWxService(db).RegisterHandler(r)
	service.NewRewardService(db).RegisterHandler(r)
	service.NewImageService(db, conf).RegisterHandler(r)
	service.NewCounterService(db).RegisterHandler(r)
	service.NewMyCourseService(db).RegisterHandler(r)
	service.NewFinishedService(db).RegisterHandler(r)

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
