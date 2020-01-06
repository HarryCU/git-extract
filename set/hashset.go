package set

import "fmt"

// see：https://www.cnblogs.com/mafeng/p/10331572.html

var Exists = struct{}{}

var StopForEachError = fmt.Errorf("StopByLoop")

type Set struct {
	bucket map[interface{}]struct{}
}

func New() *Set {
	set := &Set{bucket: make(map[interface{}]struct{})}
	return set
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.bucket[item] = Exists
	}
	return nil
}

func (s *Set) Delete(items ...interface{}) error {
	for _, item := range items {
		if s.Contains(item) {
			delete(s.bucket, item)
		}
	}
	return nil
}

func (s *Set) AddIfAbsent(items ...interface{}) error {
	for _, item := range items {
		if s.Contains(item) {
			continue
		}
		s.bucket[item] = Exists
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.bucket[item]
	return ok
}

func (s *Set) Size() int {
	return len(s.bucket)
}

func (s *Set) Clear() {
	s.bucket = make(map[interface{}]struct{})
}

func (s *Set) ForEach(callback func(interface{}) bool) (bool, error) {
	if callback == nil {
		return false, fmt.Errorf("IllegalParameter")
	}
	for key := range s.bucket {
		if !callback(key) {
			return false, StopForEachError
		}
	}
	return true, nil
}

func (s *Set) Equal(other *Set) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}

	// 迭代查询遍历
	for key := range s.bucket {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (s *Set) IsSubset(other *Set) bool {
	// s的size长于other，不用说了
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.bucket {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}
