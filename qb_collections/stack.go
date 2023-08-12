package qb_collections

import "sync"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	items []interface{}
	count int
	mux   sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewStack(initialSize int) *Stack {
	if initialSize < 0 {
		initialSize = 0
	}
	instance := new(Stack)
	instance.items = make([]interface{}, initialSize)
	instance.count = 0

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Stack) Count() int {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	return instance.count
}

// Push adds an item to the stack.
func (instance *Stack) Push(value interface{}) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	instance.items = append(instance.items[:instance.count], value)
	instance.count++
}

// Pop removes and returns an item from the stack in last to first order.
func (instance *Stack) Pop() interface{} {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if instance.count == 0 {
		return nil
	}
	instance.count--
	return instance.items[instance.count]
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
