package datasets

import (
	"sort"
	"time"
)

type ordered interface {
	time.Time | string
}

type Indexer[T ordered] struct {
	values  []T
	indices map[T]int
	inOrder bool
}

func MakeIndexer[T ordered]() *Indexer[T] {
	return &Indexer[T]{
		values:  make([]T, 0),
		indices: make(map[T]int),
		inOrder: true,
	}
}

func (idx Indexer[T]) GetIndex(value T) (index int, found bool) {
	index, found = idx.indices[value]
	return
}

func (idx Indexer[T]) Count() int {
	return len(idx.values)
}

func (idx Indexer[T]) List() (values []T) {
	if idx.inOrder == false {
		sort.Slice(idx.values, func(i, j int) bool { return isLessThan(idx.values[i], idx.values[j]) })
		idx.inOrder = true
	}
	return idx.values
}

func (idx *Indexer[T]) Add(value T) (index int, added bool) {
	var found bool
	index, found = idx.indices[value]

	if found {
		return index, false
	}

	added = true
	index = len(idx.values)
	idx.indices[value] = index

	if idx.inOrder && index > 0 {
		idx.inOrder = !isLessThan(value, idx.values[index-1])
	}
	idx.values = append(idx.values, value)
	return
}

func (idx Indexer[T]) Copy() (clone *Indexer[T]) {
	clone = &Indexer[T]{
		values:  make([]T, len(idx.values)),
		indices: make(map[T]int),
	}
	copy(clone.values, idx.values)
	for key, val := range idx.indices {
		clone.indices[key] = val
	}
	return
}

func isLessThan[T ordered](a, b T) (isLess bool) {
	// this works around the fact that we can't type switch on T
	var x interface{} = a
	var y interface{} = b
	switch (x).(type) {
	case string:
		isLess = x.(string) < y.(string)
	case time.Time:
		isLess = x.(time.Time).Before(y.(time.Time))
	}
	return
}
