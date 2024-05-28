package model

type EntityModel[K any] interface {
	// PK get primary key
	PK() K
	// SetPK set primary key
	SetPK(k K)
	// TName returns table name
	TName() string
	// Create an empty instance
	CreateInstance() any
	// ForEach iterates through fields and call given callback
	ForEach(targetValue any, sharedData any, callback func(fieldName string, fieldValue any, shared any))
}
