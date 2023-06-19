package models

import (
	"gorm.io/gorm"
)

type UserList struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);" json:"identity"` // 用户的唯一标识
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`
	AvatarUrl string `gorm:"column:avatar_url;type:varchar(100);" json:"avatar_url"`
	UserName  string `gorm:"column:user_name;type:varchar(100);" json:"user_name"`
	Password  string `gorm:"column:password;type:varchar(32);" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Mail      string `gorm:"column:mail;type:varchar(100);" json:"mail"`
	FinishNum int64  `gorm:"column:finish_num;type:int(11);" json:"finish_num"`
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"`
}

// 获取个人信息
func GetUserInfo(identity string) (info *UserList, err error) {
	info = new(UserList)
	err = DB.Where("identity = ?", identity).First(&info).Error
	if err != nil {
		return nil, err
	}
	return
}

func (table *UserList) TableName() string {
	return "user_list"
}
