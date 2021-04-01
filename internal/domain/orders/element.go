package orders

import "github.com/insan1k/proto-order-data/internal/domain/order"

//Tag describe non-trivial information about an Element which
// essentially means information about a particular element in relation to other elements
type Tag int

const (
	Empty Tag = iota
	Newest
	Matched
	Oldest
	Open
	Biggest
	Smallest
	Lowest
	Highest
	First
	Last
)

//Element this is the element as it will be used by the ring it wraps the OrderType with an index
type Element struct {
	Index int
	Order *order.Order
	tags  *[]Tag
}

// todo: transfer this logic into convenience functions - IsOldest IsNewest SetNewest SetOldest
//  ideally the type should answer a few questions about it self as
//  opposed us checking things in other packages, this is specially
//  useful in implementing logic that relates to orders, as you can
//  have the orders hold information about itself that is relevant
//  for a group of orders, i.e: is this the newest order?
//SetTags set one or more tags
func (i *Element) SetTags(o ...Tag) {
	if i.tags == nil {
		var tags []Tag
		tags = append(tags, o...)
		i.tags = &tags
	} else {
		*i.tags = append(*i.tags, o...)
	}
}

//CheckTags check one or more tags exists
func (i Element) CheckTags(o ...Tag) bool {
	if len(o) == 0 {
		if i.tags == nil {
			return true
		}
		if len(*i.tags) == 0 {
			return true
		}
	}
	for _, tag := range o {
		if !i.CheckTag(tag) {
			return false
		}
	}
	return true
}

//RemoveTag remove a particular tag
func (i *Element) RemoveTag(t Tag) bool {
	if i.tags == nil {
		return false
	}
	found := false
	var newTags []Tag
	for _, tag := range *i.tags {
		if tag != t {
			newTags = append(newTags, tag)
		} else {
			found = true
		}
	}
	i.tags = &newTags
	return found
}

//CheckTag check if a particular tag exists
func (i Element) CheckTag(o Tag) bool {
	if i.tags == nil {
		return false
	}
	for _, tag := range *i.tags {
		if tag == o {
			return true
		}
	}
	return false
}
