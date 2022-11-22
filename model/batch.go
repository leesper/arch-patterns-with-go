package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

var (
	EtaNone = time.Time{}
)

type Batch struct {
	Reference         string
	Sku               string
	PurchasedQuantity int
	Eta               time.Time
	Allocations       map[OrderLine]bool
}

func NewBatch(ref, sku string, qty int, eta time.Time) *Batch {
	return &Batch{ref, sku, qty, eta, map[OrderLine]bool{}}
}

func (b *Batch) Allocate(line OrderLine) {
	if b.CanAllocate(line) {
		b.Allocations[line] = true
	}
}

func (b *Batch) Deallocate(line OrderLine) {
	delete(b.Allocations, line)
}

func (b Batch) AvailableQuantity() int {
	return b.PurchasedQuantity - b.allocatedQuantity()
}

func (b Batch) CanAllocate(line OrderLine) bool {
	return b.Sku == line.Sku && b.AvailableQuantity() >= line.Quantity
}

func (b Batch) allocatedQuantity() int {
	allocated := 0
	for line := range b.Allocations {
		allocated += line.Quantity
	}
	return allocated
}

func (b Batch) Equals(other interface{}) bool {
	otherBatch, ok := other.(*Batch)
	if !ok {
		return ok
	}
	return otherBatch.Reference == b.Reference
}

func (b Batch) Hash() string {
	h := sha256.New()
	h.Write([]byte(b.Reference))
	return hex.EncodeToString(h.Sum(nil))
}
