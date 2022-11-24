package main

import (
	"testing"

	"github.com/leesper/arch-patterns-with-go/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	err = db.AutoMigrate(&Batch{}, &OrderLine{})
	if err != nil {
		panic("failed to create tables: " + err.Error())
	}
}
func TestRetrievingOrderLines(t *testing.T) {
	db.Exec(`INSERT INTO order_lines (order_id, sku, quantity)
		VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?)`,
		"order1", "RED-CHAIR", 12,
		"order1", "RED-TABLE", 13,
		"order2", "BLUE-LIPSTICK", 14)

	expected := []OrderLine{
		{ID: 1, BatchID: 0, OrderLine: model.OrderLine{OrderID: "order1", Sku: "RED-CHAIR", Quantity: 12}},
		{ID: 2, BatchID: 0, OrderLine: model.OrderLine{OrderID: "order1", Sku: "RED-TABLE", Quantity: 13}},
		{ID: 3, BatchID: 0, OrderLine: model.OrderLine{OrderID: "order2", Sku: "BLUE-LIPSTICK", Quantity: 14}},
	}

	var all []OrderLine
	db.Find(&all)

	if !orderLinesEqual(expected, all) {
		t.Fatalf("db.Find() == %v, want %v", all, expected)
	}
}

// TODO: 修改成泛型函数
func orderLinesEqual(expected, actual []OrderLine) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i, line := range expected {
		if line != actual[i] {
			return false
		}
	}

	return true
}

// TODO: FIXME
func TestSavingOrderLines(t *testing.T)  {}
func TestRetrievingBatches(t *testing.T) {}
func TestSavingBatches(t *testing.T)     {}
