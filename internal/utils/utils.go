package utils

import (
	"fmt"
	"sort"
)

type OrderedPair interface {
	First() string
	Last() string
}

type orderedPair struct {
	first, last string
}

func (op orderedPair) First() string {
	return op.first
}

func (op orderedPair) Last() string {
	return op.last
}

func NewOrderedPair(first, last string) OrderedPair {
	return orderedPair{
		first: first,
		last:  last,
	}
}

func Overlaps(items []OrderedPair) error {
	sort.Slice(items, func(i, j int) bool {
		return items[i].First() < items[j].First()
	})

	for i := 1; i < len(items); i++ {
		if l, f := items[i-1].Last(), items[i].First(); l >= f {
			if l == f {
				return fmt.Errorf("The element %s is supplied by more than one range.", f)
			}
			return fmt.Errorf("The elements between %s and %s are supplied by more than one range.", f, l)
		}
	}

	return nil
}
