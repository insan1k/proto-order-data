package order

//Tags describe information about an order
type Tags int

const (
	Empty Tags = iota
	Open
	Matched
	Closed
)