package main

import (
	"testing"
)

func makeBatchAndOrderLine(sku string, batchQuantity, lineQuantity int) (*Batch, OrderLine) {
	return NewBatch("batch1", sku, batchQuantity, etaNone), OrderLine{"order1", sku, lineQuantity}
}

func TestAllocatingToABatchReducesAvailableQuantity(t *testing.T) {
	batch, orderLine := makeBatchAndOrderLine("SMALL-TABLE", 20, 2)
	batch.Allocate(orderLine)
	if batch.AvailableQuantity() != 18 {
		t.Fatalf("batch.AvailableQuantity() == %d, want %d", batch.AvailableQuantity(), 18)
	}
}

func TestCanAllocateIfAvailableGreaterThanRequired(t *testing.T) {
	batch, orderLine := makeBatchAndOrderLine("BLUE-CUSHION", 10, 2)
	if !batch.CanAllocate(orderLine) {
		t.Fatalf("batch.CanAllocate(%v) == %t, want %t", orderLine, false, true)
	}
}

func TestCanAllocateIfAvailableEqualToRequired(t *testing.T) {
	batch, orderLine := makeBatchAndOrderLine("BLUE-CUSHION", 2, 2)
	if !batch.CanAllocate(orderLine) {
		t.Fatalf("batch.CanAllocate(%v) == %t, want %t", orderLine, false, true)
	}
}
func TestCannotAllocateIfAvailableLessThanRequired(t *testing.T) {
	batch, orderLine := makeBatchAndOrderLine("BLUE-CUSHION", 2, 3)
	if batch.CanAllocate(orderLine) {
		t.Fatalf("batch.CanAllocate(%v) == %t, want %t", orderLine, true, false)
	}
}

func TestCannotAllocateIfSkuNotMatch(t *testing.T) {
	batch := NewBatch("batch1", "BLUE-VASE", 10, etaNone)
	diffOrderLine := OrderLine{"order123", "BLUE-CUSHION", 2}
	if batch.CanAllocate(diffOrderLine) {
		t.Fatalf("batch.CanAllocate(%v) == %t, want %t", diffOrderLine, true, false)
	}
}

func TestAllocateToBatchShouldBeIdempotent(t *testing.T) {
	batch, orderLine := makeBatchAndOrderLine("BLUE-VASE", 10, 2)
	batch.Allocate(orderLine)
	batch.Allocate(orderLine)
	if batch.AvailableQuantity() != 8 {
		t.Fatalf("batch.AvailableQuantity() == %d, want %d", batch.AvailableQuantity(), 8)
	}
}

func TestCanOnlyDeallocateAllocatedLines(t *testing.T) {
	batch, unallocated := makeBatchAndOrderLine("BLUE-VASE", 10, 3)
	batch.Deallocate(unallocated)
	if batch.AvailableQuantity() != 10 {
		t.Fatalf("batch.AvailableQuantity() == %d, want %d", batch.AvailableQuantity(), 10)
	}
}

func TestBatchEqualityCanOnlyBasedOnReference(t *testing.T) {
	batch1 := NewBatch("batch1", "BLUE-VASE", 10, etaNone)
	batch2 := NewBatch("batch1", "BLUE-CUSHION", 2, etaNone)

	if !batch1.Equals(batch2) {
		t.Fatalf("batch1.Equals(batch2) == %t, want %t", batch1.Equals(batch2), true)
	}
}

func TestBatchHashEqualToAnotherWithSameReference(t *testing.T) {
	batch1 := NewBatch("batch1", "BLUE-VASE", 10, etaNone)
	batch2 := NewBatch("batch1", "BLUE-CUSHION", 2, etaNone)
	if batch1.Hash() != batch2.Hash() {
		t.Fatalf("batch1.Hash() == %s, want %s", batch1.Hash(), batch2.Hash())
	}
}
