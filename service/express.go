package service

import (
	"express-service/define"
	"express-service/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// @Tags 快递订单
// @Summary 订单列表
// @Description 这是一个订单列表接口
// @Router /express/list [get]
// @Param page query string true "页码"
// @Param size query string true "条数"
// @Param status query int false "接单状态"
// @Param receiver query string false "接单状态"
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
	receiver := c.Query("receiver")
	page = (page - 1) * size
	var count int64

	list := make([]*models.ExpressList, 0)
	tx := models.GetExpressList(status, receiver)
	err = tx.Omit("content").Offset(page).Limit(size).Find(&list).Count(&count).Error
	if err != nil {
		fmt.Printf("Get Express Error:", err)
		return
	}

	fmt.Println(list, "111")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
