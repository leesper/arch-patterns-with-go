package model

import (
	"fmt"
	"sort"
)

func Allocate(batches []*Batch, line OrderLine) (string, error) {
	// sort batches by eta
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].Eta.Before(batches[j].Eta)
	})

	var ref string
	var found bool
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			ref = batch.Reference
			found = true
			break
		}
	}
	if !found {
		return "", &OutOfStockErr{line.Sku}
	}
	return ref, nil
}

type OutOfStockErr struct {
	Sku string
}

func (e *OutOfStockErr) Error() string {
	return fmt.Sprintf("Out of stock for sku %s", e.Sku)
}
