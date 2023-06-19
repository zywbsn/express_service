package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// 获取 uuid
func GetUUID() string {
	s := uuid.NewV4().String()
	return s
}

// 发送验证码
func SentCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <ooooooooooos@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("<h1>您的验证码是：</h1><p>" + code + "</p>")
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "ooooooooooos@163.com", "CHUZMVVAHNSXVBCL", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	return err
}

// 生成 md5
func GetMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 生成 token
func GenerateToken(identity string) (string, error) {
	userClaims := UserClaims{
		Identity:       identity,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析 token
func AnalyseToken(token string) (*UserClaims, error) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.4inO9HZINmKFYO9qEF2SYYPHk0GuuA-qUdwIhUa8USE"
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("Analyse Token Error:%v\n", err)
	}
	return userClaims, nil
}
