package models

import (
	"gorm.io/gorm"
)

type ExpressList struct {
	gorm.Model
	CreateId      string `gorm:"column:create_id;type:varchar(255);" json:"create_id"`           // 创建人 id
	CreatePhone   string `gorm:"column:create_phone;type:varchar(255);" json:"create_phone"`     // 创建人手机号码
	Code          string `gorm:"column:code;type:varchar(255);" json:"code"`                     // 快递单号
	Address       string `gorm:"column:address;type:varchar(255);" json:"address"`               // 收货地址
	ReceiveDate   string `gorm:"column:receive_date;type:varchar(255);" json:"receive_date"`     // 收货日期
	Price         int    `gorm:"column:price;type:int(11);" json:"price"`                        // 订单费用
	ReceiveCode   string `gorm:"column:receive_code;type:varchar(255);" json:"receive_code"`     // 取件码
	Status        int    `gorm:"column:status;type:tinyint(1);" json:"status"`                   // 是否接单 0 - 否 1 - 是
	ReceiverId    string `gorm:"column:receiver_id;type:varchar(255);" json:"receiver_id"`       // 接单人 id
	ReceiverPhone string `gorm:"column:receiver_phone;type:varchar(255);" json:"receiver_phone"` // 接单人手机号码
	OrderStatus   int    `gorm:"column:order_status;type:tinyint(1);" json:"order_status"`       // 订单状态 0 - 未接单  1 - 已完成  2 - 已接单  3 - 已收货
}

type ReturnExpressList struct {
	*ExpressList
	CreateBy  string `json:"create_by"`  // 创建人姓名
	CreateImg string `json:"create_img"` // 创建人头像
	Receiver  string `json:"receiver"`   // 接单人姓名
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
