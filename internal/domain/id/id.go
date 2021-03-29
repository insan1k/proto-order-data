package id

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"time"
)

// ID defines how we're going to handle IDs in this system
type ID struct {
	uuid uuid.UUID
}

// NewID creates and return a new v4 UUID we would like the IDs to have a timestamp yes
func NewID() (i ID, err error) {
	if id, err := uuid.NewUUID(); err == nil {
		i.uuid = id
	}
	return
}

//Time returns the time form of the UUIDv1
func (i ID) Time() time.Time {
	nsec := int64(i.uuid.Time() * 100)
	return time.Unix(0, nsec)
}

//String returns the string form of the UUID
func (i ID) String() string {
	return i.uuid.String()
}

//Encode encodes the ID as base58
func (i ID) Encode() (string, error) {
	b, err := i.uuid.MarshalBinary()
	if err != nil {
		return "", err
	}
	encoded := base58.Encode(b)
	return encoded, nil
}

//Decode decodes the ShortURL.ShortID from base58
func (i *ID) Decode(b58 string) (err error) {
	var uID uuid.UUID
	err = uID.UnmarshalBinary(base58.Decode(b58))
	i.uuid = uID
	return
}
