package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

var (
	etaNone = time.Time{}
)

type Batch struct {
	reference         string
	sku               string
	purchasedQuantity int
	eta               time.Time
	allocations       map[OrderLine]bool
}

func NewBatch(ref, sku string, qty int, eta time.Time) *Batch {
	return &Batch{ref, sku, qty, eta, map[OrderLine]bool{}}
}

func (b *Batch) Allocate(line OrderLine) {
	if b.CanAllocate(line) {
		b.allocations[line] = true
	}
}

func (b *Batch) Deallocate(line OrderLine) {
	delete(b.allocations, line)
}

func (b Batch) AvailableQuantity() int {
	return b.purchasedQuantity - b.allocatedQuantity()
}

func (b Batch) CanAllocate(line OrderLine) bool {
	return b.sku == line.sku && b.AvailableQuantity() >= line.quantity
}

func (b Batch) allocatedQuantity() int {
	allocated := 0
	for line := range b.allocations {
		allocated += line.quantity
	}
	return allocated
}

func (b Batch) Equals(other interface{}) bool {
	otherBatch, ok := other.(*Batch)
	if !ok {
		return ok
	}
	return otherBatch.reference == b.reference
}

func (b Batch) Hash() string {
	h := sha256.New()
	h.Write([]byte(b.reference))
	return hex.EncodeToString(h.Sum(nil))
}

type OrderLine struct {
	orderID  string
	sku      string
	quantity int
}

func allocate(batches []*Batch, line OrderLine) (string, error) {
	// sort batches by eta
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].eta.Before(batches[j].eta)
	})

	var ref string
	var found bool
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			ref = batch.reference
			found = true
			break
		}
	}
	if !found {
		return "", &OutOfStockErr{line.sku}
	}
	return ref, nil
}

type OutOfStockErr struct {
	sku string
}

func (e *OutOfStockErr) Error() string {
	return fmt.Sprintf("Out of stock for sku %s", e.sku)
}
