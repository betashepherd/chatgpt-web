package pq

import (
	"container/heap"
	"sync"
)

// Base on Golang official demo
// An Item is something we manage in a priority queue.
type Item struct {
	Data     interface{} // The value of the item; arbitrary.
	Priority int         // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	Index int // The index of the item in the heap.
}

// A SafePriorityQueue implements heap.Interface and holds Items.
type SafePriorityQueue struct {
	IsMax bool // set type: max or min heap
	Lock  sync.RWMutex
	Items []*Item
}

func (spq SafePriorityQueue) Len() int {
	spq.Lock.RLock()
	size := len(spq.Items)
	spq.Lock.RUnlock()
	return size
}

func (spq SafePriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	if !spq.IsMax {
		return spq.Items[i].Priority < spq.Items[j].Priority
	}
	return spq.Items[i].Priority > spq.Items[j].Priority
}

func (spq SafePriorityQueue) Swap(i, j int) {
	spq.Items[i], spq.Items[j] = spq.Items[j], spq.Items[i]
	spq.Items[i].Index = i
	spq.Items[j].Index = j
}

func (spq *SafePriorityQueue) Push(x interface{}) {
	spq.Lock.Lock()
	defer spq.Lock.Unlock()
	n := len(spq.Items)
	item := x.(*Item)
	item.Index = n
	spq.Items = append(spq.Items, item)
}

func (spq *SafePriorityQueue) Pop() interface{} {
	spq.Lock.Lock()
	defer spq.Lock.Unlock()
	old := *spq
	n := len(old.Items)
	item := old.Items[n-1]
	old.Items[n-1] = nil
	spq.Items = old.Items[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (spq *SafePriorityQueue) Update(item *Item, priority int) {
	item.Priority = priority
	heap.Fix(spq, item.Index)
}
