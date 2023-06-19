package main

// @title 快递代取 API
// @version 1.0.1
// @contact.name silenceLamb
// @contact.url http://www.swagger.io/support
// @contact.email ooooooooooos@163.com
// @host localhost:9090
// @BasePath /

import (
	_ "express-service/docs"
	"express-service/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := router.Router()
	r.Static("/static/images", "./static/images")
	// 注册Swagger接口文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":9090")
}
