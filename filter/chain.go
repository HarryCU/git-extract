package filter

import (
	"github.com/HarryCU/git-extract/set"
)

type Chain struct {
	items *set.Set
}

func New(filters ...Filter) *Chain {
	chain := &Chain{
		items: set.New(),
	}
	if len(filters) != 0 {
		for _, filter := range filters {
			_ = chain.items.AddIfAbsent(filter)
		}
	}
	return chain
}

func (fc *Chain) Include(value interface{}) bool {
	if fc.items.Size() == 0 {
		return false
	}
	r, _ := fc.items.ForEach(func(v interface{}) bool {
		filter := v.(Filter)
		if !filter.Include(value) {
			return false
		}
		return true
	})
	return r
}
