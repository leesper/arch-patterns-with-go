package main

import (
	"testing"
	"time"

	"github.com/leesper/arch-patterns-with-go/model"
)

func TestPreferBatchesInWarehouseRatherThanShipments(t *testing.T) {
	batchInStock := model.NewBatch("batch1", "BLUE-VASE", 10, model.EtaNone)
	now := time.Now()
	batchOnShip := model.NewBatch("batch2", "BLUE-VASE", 10, now)
	orderLine := model.OrderLine{OrderID: "oref", Sku: "BLUE-VASE", Quantity: 2}
	model.Allocate([]*model.Batch{batchInStock, batchOnShip}, orderLine)
	if batchInStock.AvailableQuantity() != 8 {
		t.Fatalf("batchInStock.AvailableQuantity() == %d, want %d", batchInStock.AvailableQuantity(), 8)
	}
	if batchOnShip.AvailableQuantity() != 10 {
		t.Fatalf("batchOnShip.AvailableQuantity() == %d, want %d", batchOnShip.AvailableQuantity(), 10)
	}
}
func TestPreferEarliestBatches(t *testing.T) {
	current := time.Now()
	theDayBefore := current.Add(time.Hour * -24)
	theDayAfter := current.Add(time.Hour * 24)
	batch1 := model.NewBatch("batch1", "BLUE-VASE", 10, theDayBefore)
	batch2 := model.NewBatch("batch2", "BLUE-VASE", 10, current)
	batch3 := model.NewBatch("batch3", "BLUE-VASE", 10, theDayAfter)
	orderLine := model.OrderLine{OrderID: "oref", Sku: "BLUE-VASE", Quantity: 2}
	model.Allocate([]*model.Batch{batch1, batch2, batch3}, orderLine)

	if batch1.AvailableQuantity() != 8 {
		t.Fatalf("batch1.AvailableQuantity() == %d, want %d", batch1.AvailableQuantity(), 8)
	}

	if batch2.AvailableQuantity() != 10 {
		t.Fatalf("batch2.AvailableQuantity() == %d, want %d", batch2.AvailableQuantity(), 10)
	}

	if batch3.AvailableQuantity() != 10 {
		t.Fatalf("batch3.AvailableQuantity() == %d, want %d", batch3.AvailableQuantity(), 10)
	}
}

func TestReturnsAllocatedBatchRef(t *testing.T) {
	batchInStock := model.NewBatch("batch1", "BLUE-VASE", 10, model.EtaNone)
	batchOnShip := model.NewBatch("batch2", "BLUE-VASE", 10, time.Now())
	orderLine := model.OrderLine{OrderID: "oref", Sku: "BLUE-VASE", Quantity: 2}
	batchRef, _ := model.Allocate([]*model.Batch{batchInStock, batchOnShip}, orderLine)
	if batchInStock.Reference != batchRef {
		t.Fatalf("batchInStock.reference == %s, want %s", batchInStock.Reference, batchRef)
	}
}

func TestReturnOutOfStockErrIfCannotAllocate(t *testing.T) {
	batch := model.NewBatch("batch2", "BLUE-VASE", 10, time.Now())
	line := model.OrderLine{OrderID: "oref", Sku: "BLUE-VASE", Quantity: 10}
	model.Allocate([]*model.Batch{batch}, line)

	_, err := model.Allocate([]*model.Batch{batch}, model.OrderLine{OrderID: "oref", Sku: "BLUE-VASE", Quantity: 2})
	e, ok := err.(*model.OutOfStockErr)
	if !ok {
		t.Fatalf("allocate() returns error %T, want %T", e, &model.OutOfStockErr{})
	}
	if e.Sku != "BLUE-VASE" {
		t.Fatalf("e.sku == %s, want %s", e.Sku, "BLUE-VASE")
	}
}
