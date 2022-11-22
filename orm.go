package main

import (
	"time"

	"github.com/leesper/arch-patterns-with-go/model"
)

// TODO: 创建batch和orderline的数据库记录表，并建立它们之间one-to-many的关系
type BatchRecord struct {
	ID                uint
	Reference         string
	Sku               string
	PurchasedQuantity int
	Eta               time.Time
	Allocations       []model.OrderLine
}

func (BatchRecord) TableName() string {
	return "batches"
}

type OrderLineRecord struct {
	ID            uint
	BatchRecordID uint
	model.OrderLine
}

func (OrderLineRecord) TableName() string {
	return "order_lines"
}
