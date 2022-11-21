package main

import (
	"testing"
	"time"
)

func TestPreferBatchesInWarehouseRatherThanShipments(t *testing.T) {
	batchInStock := NewBatch("batch1", "BLUE-VASE", 10, etaNone)
	now := time.Now()
	batchOnShip := NewBatch("batch2", "BLUE-VASE", 10, now)
	orderLine := OrderLine{"oref", "BLUE-VASE", 2}
	allocate([]*Batch{batchInStock, batchOnShip}, orderLine)
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
	batch1 := NewBatch("batch1", "BLUE-VASE", 10, theDayBefore)
	batch2 := NewBatch("batch2", "BLUE-VASE", 10, current)
	batch3 := NewBatch("batch3", "BLUE-VASE", 10, theDayAfter)
	orderLine := OrderLine{"oref", "BLUE-VASE", 2}
	allocate([]*Batch{batch1, batch2, batch3}, orderLine)

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
	batchInStock := NewBatch("batch1", "BLUE-VASE", 10, etaNone)
	batchOnShip := NewBatch("batch2", "BLUE-VASE", 10, time.Now())
	orderLine := OrderLine{"oref", "BLUE-VASE", 2}
	batchRef, _ := allocate([]*Batch{batchInStock, batchOnShip}, orderLine)
	if batchInStock.reference != batchRef {
		t.Fatalf("batchInStock.reference == %s, want %s", batchInStock.reference, batchRef)
	}
}

func TestReturnOutOfStockErrIfCannotAllocate(t *testing.T) {
	batch := NewBatch("batch2", "BLUE-VASE", 10, time.Now())
	line := OrderLine{"oref", "BLUE-VASE", 10}
	allocate([]*Batch{batch}, line)

	_, err := allocate([]*Batch{batch}, OrderLine{"oref", "BLUE-VASE", 2})
	e, ok := err.(*OutOfStockErr)
	if !ok {
		t.Fatalf("allocate() returns error %T, want %T", e, &OutOfStockErr{})
	}
	if e.sku != "BLUE-VASE" {
		t.Fatalf("e.sku == %s, want %s", e.sku, "BLUE-VASE")
	}
}
