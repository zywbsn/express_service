package models

import (
	"gorm.io/gorm"
)

type ExpressList struct {
	gorm.Model
	CreateBy    string `gorm:"column:create_by;type:varchar(255);" json:"create_by"`       // 创建人姓名
	Code        string `gorm:"column:code;type:varchar(255);" json:"code"`                 // 快递单号
	Address     string `gorm:"column:address;type:varchar(255);" json:"address"`           // 收货地址
	ReceiveDate string `gorm:"column:receive_date;type:varchar(255);" json:"receive_date"` // 收货日期
	Price       int    `gorm:"column:price;type:int(11);" json:"price"`                    // 订单费用
	ReceiveCode string `gorm:"column:receive_code;type:varchar(255);" json:"receive_code"` // 取件码
	Status      int    `gorm:"column:status;type:tinyint(1);" json:"status"`               // 是否接单 0 - 否 1 - 是
	Receiver    string `gorm:"column:receiver;type:varchar(255);" json:"receiver"`         // 接单人姓名
	ReceiverId  string `gorm:"column:receiver_id;type:varchar(255);" json:"receiver_id"`   // 接单人 id
	CreateId    string `gorm:"column:create_id;type:varchar(255);" json:"create_id"`       // 创建人 id
	CreateImg   string `gorm:"column:create_img;type:varchar(255);" json:"create_img"`     // 创建人头像
}

func GetExpressDetail(id string) *gorm.DB {
	tx := DB.Where("id = ?", id)
	return tx
}

// 查询列表
func GetExpressList(status int, receiverId, createId string) *gorm.DB {
	tx := DB.Model(new(ExpressList))
	if status != -1 {
		tx.Where("status = ?", status)
	}
	if receiverId != "" {
		tx.Where("receiver_id = ?", receiverId)
	}
	if createId != "" {
		tx.Where("create_id = ?", createId)
	}
	return tx
}

func (table *ExpressList) TableName() string {
	return "express_list"
}
