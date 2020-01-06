package filter

type Filter interface {
	Include(value interface{}) bool
}
