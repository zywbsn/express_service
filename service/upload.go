package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

// @Tags 上传文件
// @Summary 上传文件
// @Description 上传文件接口
// @Router /upload [post]
// @Param file formData  string true "上传的文件"
// @Produce multipart/form-data
// @Success 200 {string} string
func UploadFile(context *gin.Context) {
	// 从请求中读取文件
	f, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		// 将读取的文件存储在本地（服务端本地）
		// dst := fmt.Sprint("./%s", f.Filename)
		dst := path.Join("./static/images", f.Filename)
		fmt.Println("/"+dst, "dst")
		context.SaveUploadedFile(f, dst)
		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": map[string]interface{}{
				"path":    "/" + dst,
				"message": "上传成功",
			},
		})
	}
}
