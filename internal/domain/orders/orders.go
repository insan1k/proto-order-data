package orders

import (
	"container/ring"
	"errors"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"time"
)

const (
	initialOrdersTooBigErr = "cannot allocate amount of orders greater then ring capacity, " +
		"allocate an empty Orders and populate it with another operation, or increase size of ring"
	emptyOrders = "orders struct is currently empty"
)

// Orders implements a type that holds multiple orders data... this could be the bids of an orders
// book or the X latest orders this structure implies that orders are sorted in a useful manner
type Orders struct {
	asset       string
	ring        *ring.Ring
	first       *Element
	last        *Element
	maxSize     int
	currentSize int
}

// Asset returns the asset of these current orders
func (o Orders) Asset() string {
	return o.asset
}

//NewOrders receives a list of orders and returns them in the Orders struct
func NewOrders(asset string, maxSize int, orders ...*order.Order) (o Orders, err error) {
	o.asset = asset
	if len(orders) > maxSize {
		err = errors.New(initialOrdersTooBigErr)
		return
	}
	o.ring = ring.New(maxSize)
	o.maxSize = maxSize
	isFirst := true
	isLast := false
	for index := 0; index < maxSize; index++ {
		o.ring.Value = &Element{
			Index: index,
			Order: nil,
			tags:  nil,
		}
		o.next()
	}
	o.moveToFirst()
	for c, currentOrder := range orders {
		if c == len(orders)-1 {
			isLast = true
		}
		if isFirst {
			o.first = o.get()
			o.first.SetTags(First)
			isFirst = false
		}
		if isLast {
			o.last = o.get()
			o.last.SetTags(Last)
			isLast = false
		}
		o.currentSize = c
		o.set(currentOrder)
		o.ring.Next()
	}
	return
}

// First returns a copy of the first order
func (o Orders) First() (order.Order, error) {
	if o.currentSize == 0 {
		return order.EmptyOrder(), errors.New(emptyOrders)
	}
	return *o.first.Order, nil
}

// Last returns a copy of the last order
func (o Orders) Last() (order.Order, error) {
	if o.currentSize == 0 {
		return order.EmptyOrder(), errors.New(emptyOrders)
	}
	return *o.last.Order, nil
}

// TimeStart returns the time we received the first order
func (o Orders) TimeStart() time.Time {
	return o.first.Order.Inf.Seen()
}

// TimeEnd returns the time we received the last order
func (o Orders) TimeEnd() time.Time {
	return o.last.Order.Inf.Seen()
}

// TimePeriod returns the time period of orders inside
func (o Orders) TimePeriod() time.Duration {
	return o.TimeEnd().Sub(o.TimeStart())
}

// Cap returns the maximum capacity of this Orders struct
<<<<<<< HEAD
func (o Orders) Cap() (i int) {
=======
func (o Orders) Cap()(i int){
>>>>>>> 27d5f2ff5f8f7d768344c848b2ce50316e28c857
	return o.maxSize
}

// Len returns the utilized capacity of this Orders struct
func (o Orders) Len() (i int) {
	return o.currentSize
}

func (o *Orders) next() {
	o.ring = o.ring.Next()
}

func (o *Orders) get() *Element {
	return (o.ring.Value).(*Element)
}

func (o *Orders) set(order *order.Order) {
	this := o.get()
	this.Order = order
	// notice that when setting the element we remove it's tags this is due to the fact that the tags of an element
	// may no longer make sense after setting it
	this.tags = nil
}

//Insert adds order.Order into ring
func (o *Orders) Insert(order *order.Order) {
	//check if we're operating under an empty Orders ring
	if o.currentSize == 0 {
		o.set(order)
		o.first = o.get()
		o.last = o.get()
		o.first.SetTags(First)
		o.last.SetTags(Last)
		o.currentSize++
		return
	}
	//otherwise we go straight to the last
	o.moveToLast()
	//remove it's tag as the lest element
	o.last.RemoveTag(Last)
	// the next element will be the order we're inserting, so we go there
	o.next()
	// we set it's order reference
	o.set(order)
	//it will also be the new last, so we update it's reference
	o.last = o.get()
	//and we set it's tag
	o.last.SetTags(Last)
	//if our orders are already at cap we need to move the first as well
	if o.currentSize == o.maxSize {
		// the first one will be the element immediately after the new last, so we go there
		o.next()
		// we update our reference to first
		o.first = o.get()
		// and we set it's tag
		o.first.SetTags(First)
	} else {
		//otherwise we increment the current size until we're at capacity
		o.currentSize++
	}
}

func (o *Orders) moveToLast() {
	var lastIndex int
	if o.currentSize == 0 {
		lastIndex = 0
	} else {
		lastIndex = o.last.Index
	}
	// figure out where we are in the ring
	e := o.get()
	if e.Index-lastIndex != 0 {
		//go to the last element
		o.ring = o.ring.Move(lastIndex - e.Index)
	}
}

func (o *Orders) moveToFirst() {
	o.moveToLast()
	o.next()
}

//todo: sorting functions
