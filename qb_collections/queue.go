package qb_collections

import "sync"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	items []interface{}
	size  int
	head  int
	tail  int
	count int
	mux   sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewQueue(initialSize int) *Queue {
	if initialSize < 1 {
		initialSize = 1
	}
	instance := new(Queue)
	instance.items = make([]interface{}, initialSize)
	instance.size = initialSize

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Queue) Count() int {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	return instance.count
}

// Push adds a node to the queue.
func (instance *Queue) Push(value interface{}) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if instance.head == instance.tail && instance.count > 0 {
		nodes := make([]interface{}, len(instance.items)+instance.size)
		copy(nodes, instance.items[instance.head:])
		copy(nodes[len(instance.items)-instance.head:], instance.items[:instance.head])
		instance.head = 0
		instance.tail = len(instance.items)
		instance.items = nodes
	}
	instance.items[instance.tail] = value
	instance.tail = (instance.tail + 1) % len(instance.items)
	instance.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (instance *Queue) Pop() interface{} {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if instance.count == 0 {
		return nil
	}
	item := instance.items[instance.head]
	instance.head = (instance.head + 1) % len(instance.items)
	instance.count--
	return item
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
