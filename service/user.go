package service

import (
	"encoding/json"
	"errors"
	"express-service/define"
	"express-service/helper"
	"express-service/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// @Tags 用户
// @Summary 修改用户信息
// @Description 修改用户信息接口
// @Router /user/info [put]
// @Param identity query  string true "用户唯一标识"
// @Param avatar_url query  string false "头像"
// @Param mail query  string false "电子邮箱"
// @Param name query  string false "昵称"
// @Param phone query  string false "手机号码"
// @Param user_name query  string false "账户"
// @Param password query  string false "密码"
// @Produce application/json
// @Success 200 {string} string
func UpdateUserInfo(c *gin.Context) {
	identity := c.PostForm("identity")
	avatarUrl := c.PostForm("avatar_url")
	mail := c.PostForm("mail")
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	userName := c.PostForm("user_name")
	password := c.PostForm("password")
	info := &models.UserList{
		Name:      name,
		AvatarUrl: avatarUrl,
		UserName:  userName,
		Password:  password,
		Phone:     phone,
		Mail:      mail,
	}
	err := models.DB.Model(new(models.UserList)).Where("identity = ?", identity).
		Updates(info).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "User Info err:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "修改成功",
	})
}

// @Tags 用户
// @Summary 用户详情
// @Description 用户详情接口
// @Router /user/info [get]
// @Param identity query  string true "identity"
// @Produce application/json
// @Success 200 {string} string
func GetUserInfo(c *gin.Context) {
	identity := c.Query("identity")
	info, err := models.GetUserInfo(identity)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "User Info err:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": info,
	})
}

// @Tags 用户
// @Summary 用户登录
// @Description 用户登录接口
// @Router /admin/login [post]
// @Param code formData  string true "username"
// @Param name formData  string true "username"
// @Produce application/json
// @Success 200 {string} string
func AdminLogin(c *gin.Context) {
	username := c.PostForm("user_name")
	password := c.PostForm("password")
	fmt.Println(username)
	data := new(models.UserList)
	err := models.DB.Where("user_name = ?", username).First(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"error":   err.Error(),
			"message": "该用户不存在",
		})
		return
	}
	if password != data.Password {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "账号密码错误",
		})
		return
	}
	token, _ := helper.GenerateToken(data.Identity)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登陆成功",
		"info":    data,
		"token":   token,
	})
}

// @Tags 用户
// @Summary 用户登录
// @Description 用户登录接口
// @Router /user/login [post]
// @Param code formData  string true "code"
// @Param name formData  string true "名字"
// @Param avatarUrl formData  string true "头像"
// @Produce application/json
// @Success 200 {string} string
func Login(c *gin.Context) {
	code := c.PostForm("code")            //  获取 code
	name := c.PostForm("name")            //  获取 name
	avatarUrl := c.PostForm("avatar_url") //  获取 avatarUrl

	// 根据code获取 openID 和 session_key
	wxLoginResp, err := WXLogin(code)
	fmt.Println("wxLoginResp:%v", wxLoginResp)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	token, err := helper.GenerateToken(wxLoginResp.OpenId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken Error:" + err.Error(),
		})
	}
	data := new(models.UserList)
	err = models.DB.Where("identity = ?", wxLoginResp.OpenId).First(&data).Error
	if err != nil {
		info := &models.UserList{
			Name:      name,
			AvatarUrl: avatarUrl,
			Identity:  wxLoginResp.OpenId,
		}
		err = models.DB.Create(&info).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "User Create Error:" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "登陆成功",
			"info":    info,
			"token":   token,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登陆成功",
		"info":    data,
		"token":   token,
	})
}

// 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况
func WXLogin(code string) (*WXLoginResp, error) {
	fmt.Printf("code:%v", code)
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的 appId 和 secret 是在微信公众平台上获取的
	url = fmt.Sprintf(url, define.AppId, define.Secret, code)

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	return &wxResp, nil
}
