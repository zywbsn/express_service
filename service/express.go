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
// @Param create_img formData string true "创建人头像"
// @Produce application/json
// @Success 200 {string} string
func CreateExpress(c *gin.Context) {
	Code := c.PostForm("code")
	Address := c.PostForm("address")
	ReceiveDate := c.PostForm("receive_date")
	Price, _ := strconv.Atoi(c.PostForm("price"))
	ReceiveCode := c.PostForm("receive_code")
	CreateId := c.PostForm("create_id")
	CreateImg := c.PostForm("create_img")

	userInfo := new(models.UserList)
	models.DB.Where("identity = ?", CreateId).First(&userInfo)

	data := &models.ExpressList{
		CreateBy:    userInfo.Name,
		Code:        Code,
		Address:     Address,
		ReceiveDate: ReceiveDate,
		Price:       Price,
		ReceiveCode: ReceiveCode,
		CreateId:    CreateId,
		CreateImg:   CreateImg,
	}
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "Express Create Error:" + err.Error(),
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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
