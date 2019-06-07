package data

type Paged interface {
	GetItems() interface{}
	GetTotal() int64
	GetLimit() int64
	GetOffset() int64
}
