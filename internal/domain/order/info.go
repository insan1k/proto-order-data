package order

import (
	"github.com/insan1k/proto-order-data/internal/domain/id"
	"time"
)

//Info holds information about a particular orders, from the raw data of an orders to relevant information to an orders in
//relation to other orders
type Info struct {
	id   id.ID
	meta *[]byte
	tags *[]Tags
}

func (i *Info) init() {
	newId, err := id.NewID()
	if err != nil {
		panic(err)
	}
	i.id = newId
}

//SetMeta stores metadata i.e: the raw data we got from the exchange
func (i *Info) SetMeta(b []byte) {
	i.meta = &b
}

//GetMeta retrieves a copy of the metadata
func (i Info) GetMeta() (b []byte) {
	return *i.meta
}

// todo: transfer this logic into convenience functions
//SetTags
func (i *Info) SetTags(o ...Tags) {
	if i.tags == nil {
		var tags []Tags
		tags = append(tags, o...)
		i.tags = &tags
	} else {
		*i.tags = append(*i.tags, o...)
	}
}

//CheckTags
func (i Info) CheckTags(o ...Tags) bool {
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

//RemoveTag
func (i *Info) RemoveTag(o Tags) bool {
	if i.tags == nil {
		return false
	}
	found := false
	var newTags []Tags
	for _, tag := range *i.tags {
		if tag != o {
			newTags = append(newTags, tag)
		} else {
			found = true
		}
	}
	i.tags = &newTags
	return found
}

//CheckTag
func (i Info) CheckTag(o Tags) bool {
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

//Seen the creation time of this Order Info
func (i Info) Seen() time.Time {
	return i.id.Time()
}
