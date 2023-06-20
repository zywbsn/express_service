package router

import (
	"express-service/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	/* ----- 订单相关路由组 ----- */
	expressGroup := r.Group("/express")
	{
		// 订单列表
		expressGroup.GET("/list", service.GetExpressList)   // 代取列表
		expressGroup.GET("/info", service.GetExpressDetail) // 订单详情
		expressGroup.POST("/create", service.CreateExpress) // 新建订单
	}
	/* ----- 订单相关路由组 ----- */

	/* ----- 用户相关路由组 ----- */
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", service.Login)        // 用户存在登录 不存在注册并登录
		userGroup.GET("/info", service.GetUserInfo)    // 获取用户个人信息
		userGroup.PUT("/info", service.UpdateUserInfo) // 修改用户个人信息
	}
	/* ----- 用户相关路由组 ----- */

	/* ----- 上传文件 ----- */
	r.POST("/upload", service.UploadFile) // 上传文件
	/* ----- 上传文件 ----- */

	return r
}
