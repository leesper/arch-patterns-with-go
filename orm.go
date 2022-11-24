package main

import (
	"time"

	"github.com/leesper/arch-patterns-with-go/model"
)

// TODO: 创建batch和orderline的数据库记录表，并建立它们之间one-to-many的关系
type Batch struct {
	ID                uint `gorm:"primaryKey;autoIncrement:true"`
	Reference         string
	Sku               string
	PurchasedQuantity int
	Eta               time.Time
	Allocations       []OrderLine
}

type OrderLine struct {
	ID      uint `gorm:"primaryKey;autoIncrement:true"`
	BatchID uint
	model.OrderLine
}
