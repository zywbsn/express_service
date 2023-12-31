package router

import (
	"express-service/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api")
	{
		/* ----- 后台管理 ----- */
		apiGroup.POST("/admin/login", service.AdminLogin)
		/* ----- 后台管理 ----- */
		/* ----- 订单相关路由组 ----- */
		expressGroup := apiGroup.Group("/express")
		{
			// 订单列表
			expressGroup.GET("/list", service.GetExpressList)   // 代取列表
			expressGroup.GET("/info", service.GetExpressDetail) // 订单详情
			expressGroup.POST("/create", service.CreateExpress) // 新建订单
			expressGroup.PUT("/order", service.TakeOrder)       // 接单接口
			expressGroup.PUT("/finish", service.FinishOrder)    // 完成订单接口
		}
		/* ----- 订单相关路由组 ----- */

		/* ----- 用户相关路由组 ----- */
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/login", service.Login)        // 用户存在登录 不存在注册并登录
			userGroup.GET("/info", service.GetUserInfo)    // 获取用户个人信息
			userGroup.PUT("/info", service.UpdateUserInfo) // 修改用户个人信息
		}
		/* ----- 用户相关路由组 ----- */

		/* ----- 上传文件 ----- */
		apiGroup.POST("/upload", service.UploadFile) // 上传文件
		/* ----- 上传文件 ----- */
	}

	return r
}
