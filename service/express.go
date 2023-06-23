package service

import (
	"express-service/define"
	"express-service/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// @Tags 快递订单
// @Summary 完成订单
// @Description 完成订单接口 - 传入接单人 identity 为完成订单 传入创建人 identity 为收货
// @Router /express/finish [put]
// @Param id query string true "订单 id"
// @Param identity query string true "用户唯一标识"
// @Produce application/json
// @Success 200 {string} string
func FinishOrder(c *gin.Context) {
	id := c.Query("id")
	identity := c.Query("identity")
	info := new(models.ExpressList)
	tx := models.DB.Model(new(models.ExpressList)).Where("id = ?", id)
	err := tx.First(&info).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "订单不存在",
		})
		return
	}

	// 收货逻辑
	if info.CreateId == identity {
		if info.OrderStatus != 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "该订单未完成",
			})
			return
		}
		info.OrderStatus = 3
		err = tx.Updates(info).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"error":   err.Error(),
				"message": "确认收货失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "收货成功",
		})
		return
	}

	// 完成订单逻辑
	if info.ReceiverId == identity {
		if info.OrderStatus != 2 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "参数错误",
			})
			return
		}
		info.OrderStatus = 1
		err = tx.Updates(info).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"error":   err.Error(),
				"message": "订单完成失败",
			})
			return
		}

		userInfo, _ := models.GetUserInfo(identity)
		userInfo.FinishNum++
		err = models.DB.Model(new(models.UserList)).Where("identity = ?", identity).Updates(userInfo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"error":   err.Error(),
				"message": "订单完成失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "订单完成",
		})
	}

}

// @Tags 快递订单
// @Summary 接单
// @Description 这是一个接单接口
// @Router /express/order [put]
// @Param id query string true "订单 id"
// @Param receiver_id query string true "接单人 id"
// @Param receiver_phone query string true "接单人手机号码"
// @Produce application/json
// @Success 200 {string} string
func TakeOrder(c *gin.Context) {
	id := c.Query("id")
	info := new(models.ExpressList)
	tx := models.GetExpressDetail(id)
	err := tx.First(&info).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "订单不存在",
		})
		return
	}

	if info.Status == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "该订单已被接单",
		})
		return
	}

	receiverId := c.Query("receiver_id")

	_, err = models.GetUserInfo(receiverId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "用户不存在",
		})
		return
	}
	if info.CreateId == receiverId {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "当前用户为创建人",
		})
		return
	}

	info.Status = 1
	info.ReceiverId = receiverId
	info.OrderStatus = 2
	info.ReceiverPhone = c.Query("receiver_phone")

	err = tx.Updates(info).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "接单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "接单成功",
	})
}

// @Tags 快递订单
// @Summary 订单详情
// @Description 这是一个订单详情接口
// @Router /express/info [get]
// @Param id query string true "订单 id"
// @Produce application/json
// @Success 200 {string} string
func GetExpressDetail(c *gin.Context) {
	id := c.Query("id")
	info := new(models.ExpressList)
	err := models.DB.Where("id = ?", id).First(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "订单不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "Get ExpressDetail Error：" + err.Error(),
		})
		return
	}
	info.CreateBy = models.GetName(info.CreateId)
	info.CreateImg = models.GetImage(info.CreateId)
	info.Receiver = models.GetName(info.ReceiverId)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"data": info,
		},
	})
}

// @Tags 快递订单
// @Summary 新增订单
// @Description 这是一个新增订单接口
// @Router /express/create [post]
// @Param code formData string true "快递单号"
// @Param address formData string true "收货地址"
// @Param receive_date formData string true "收货日期"
// @Param price formData string true "订单费用"
// @Param receive_code formData string true "取件码"
// @Param create_id formData string true "创建人 id"
// @Param create_phone formData string true "创建人手机号"
// @Produce application/json
// @Success 200 {string} string
func CreateExpress(c *gin.Context) {
	Code := c.PostForm("code")
	Address := c.PostForm("address")
	ReceiveDate := c.PostForm("receive_date")
	Price, _ := strconv.Atoi(c.PostForm("price"))
	ReceiveCode := c.PostForm("receive_code")
	CreateId := c.PostForm("create_id")
	CreatePhone := c.PostForm("create_phone")

	data := &models.ExpressList{
		CreateBy:    models.GetName(CreateId),
		Code:        Code,
		Address:     Address,
		ReceiveDate: ReceiveDate,
		Price:       Price,
		ReceiveCode: ReceiveCode,
		CreateId:    CreateId,
		CreateImg:   models.GetImage(CreateId),
		CreatePhone: CreatePhone,
	}
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "新增订单失败",
		})
		return
	}
	userInfo, _ := models.GetUserInfo(CreateId)
	userInfo.SubmitNum++
	err = models.DB.Model(new(models.UserList)).Where("identity = ?", CreateId).Updates(userInfo).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"error":   err.Error(),
			"message": "新增订单失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

// @Tags 快递订单
// @Summary 订单列表
// @Description 这是一个订单列表接口
// @Router /express/list [get]
// @Param page query string true "页码"
// @Param size query string true "条数"
// @Param status query int false "接单状态"
// @Param receiver_id query string false "接单人"
// @Param create_id query string false "创建人"
// @Produce application/json
// @Success 200 {string} string
func GetExpressList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Panicln("GetProblemList Page Parse Error:", err)
		return
	}
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	receiverId := c.Query("receiver_id")
	createId := c.Query("create_id")
	page = (page - 1) * size
	var count int64

	list := make([]*models.ExpressList, 0)
	tx := models.GetExpressList(status, receiverId, createId)
	err = tx.Omit("content").Offset(page).Limit(size).Find(&list).Count(&count).Error
	if err != nil {
		fmt.Printf("Get Express Error:", err)
		return
	}

	// 动态创建人和接单人姓名
	for i := 0; i < len(list); i++ {
		list[i].Receiver = models.GetName(list[i].ReceiverId)
		list[i].CreateBy = models.GetName(list[i].CreateId)
		list[i].CreateImg = models.GetImage(list[i].CreateId)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
