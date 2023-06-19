package models

import (
	"fmt"
	"gorm.io/gorm"
)

type ExpressList struct {
	gorm.Model
	CreateBy    string `gorm:"column:create_by;type:varchar(255);" json:"create_by"`       // 问题的唯一标识
	Code        string `gorm:"column:code;type:varchar(255);" json:"code"`                 // 关联问题分类表 id
	Address     string `gorm:"column:address;type:varchar(255);" json:"address"`           // 文章标题
	ReceiveDate string `gorm:"column:receive_date;type:varchar(255);" json:"receive_date"` // 文章正文
	Price       int    `gorm:"column:price;type:int(11);" json:"price"`                    // 最大运行时长
	ReceiveCode string `gorm:"column:receive_code;type:varchar(255);" json:"receive_code"` // 最大允许内存
	Status      int    `gorm:"column:status;type:tinyint(1);" json:"status"`               // 关联测试用表
	Receiver    string `gorm:"column:receiver;type:varchar(255);" json:"receiver"`         // 关联测试用表
}

// 查询列表
func GetExpressList(status int, receiver string) *gorm.DB {
	tx := DB.Model(new(ExpressList))
	if status != -1 {
		tx.Where("status = ?", status)
	}
	if receiver != "" {
		fmt.Println(131, receiver)
		tx.Where("receiver = ?", receiver)
	}
	return tx
}

func (table *ExpressList) TableName() string {
	return "express_list"
}
