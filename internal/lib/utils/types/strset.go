package types

import (
	"sync"
)

type StrSet struct {
	m     map[string]struct{}
	mu    sync.RWMutex
	block bool
}

func NewStrSet(block bool, values ...string) *StrSet {
	var s = &StrSet{}
	s.block = block
	s.m = make(map[string]struct{})
	if len(values) > 0 {
		s.Add(values...)
	}
	return s
}

func (ss *StrSet) lock() {
	if ss.block {
		ss.mu.Lock()
	}
}

func (ss *StrSet) unlock() {
	if ss.block {
		ss.mu.Unlock()
	}
}

func (ss *StrSet) rLock() {
	if ss.block {
		ss.mu.RLock()
	}
}

func (ss *StrSet) rUnlock() {
	if ss.block {
		ss.mu.RUnlock()
	}
}

func (ss *StrSet) Add(values ...string) {
	ss.lock()
	defer ss.unlock()

	for _, v := range values {
		ss.m[v] = struct{}{}
	}
}

func (ss *StrSet) Remove(values ...string) {
	ss.lock()
	defer ss.unlock()

	for _, v := range values {
		delete(ss.m, v)
	}
}

func (ss *StrSet) RemoveAll() {
	ss.lock()
	defer ss.unlock()

	for k := range ss.m {
		delete(ss.m, k)
	}
}

func (ss *StrSet) Exists(v string) bool {
	ss.rLock()
	defer ss.rUnlock()

	_, found := ss.m[v]
	return found
}

func (ss *StrSet) NotExists(v string) bool {
	return !ss.Exists(v)
}

func (ss *StrSet) Contains(values ...string) bool {
	ss.rLock()
	defer ss.rUnlock()

	for _, v := range values {
		if _, found := ss.m[v]; !found {
			return false
		}
	}
	return true
}

func (ss *StrSet) Has(val string) bool {
	ss.rLock()
	defer ss.rUnlock()

	_, has := ss.m[val]
	return has
}

func (ss *StrSet) NotHas(val string) bool {
	return !ss.Has(val)
}

func (ss *StrSet) Len() int {
	ss.rLock()
	defer ss.rUnlock()

	return ss.len()
}

func (ss *StrSet) len() int {
	return len(ss.m)
}

func (ss *StrSet) Values() []string {
	ss.rLock()
	defer ss.rUnlock()

	var vs = make([]string, 0, ss.len())
	for k := range ss.m {
		vs = append(vs, k)
	}
	return vs
}

func (ss *StrSet) Iter() <-chan string {
	var iv = make(chan string, len(ss.m))

	go func(s *StrSet) {
		s.rLock()
		for k := range ss.m {
			iv <- k
		}
		close(iv)
		s.rUnlock()
	}(ss)

	return iv
}

func (ss *StrSet) Equal(s *StrSet) bool {
	ss.rLock()
	defer ss.rUnlock()

	if ss.len() != s.Len() {
		return false
	}

	for k := range ss.m {
		if !s.Exists(k) {
			return false
		}
	}

	return true
}

func (ss *StrSet) Clone() *StrSet {
	ss.rLock()
	defer ss.rUnlock()

	var ns = NewStrSet(ss.block)
	for k := range ss.m {
		ns.Add(k)
	}
	return ns
}

func (ss *StrSet) Intersect(s *StrSet) *StrSet {
	ss.rLock()
	defer ss.rUnlock()

	var ns = NewStrSet(ss.block)
	var vs = s.Values()
	for _, v := range vs {
		_, exists := ss.m[v]
		if exists {
			ns.Add(v)
		}
	}
	return ns
}

func (ss *StrSet) Union(s *StrSet) *StrSet {
	ss.rLock()
	defer ss.rUnlock()

	var ns = NewStrSet(ss.block)
	ns.Add(ss.Values()...)
	ns.Add(s.Values()...)
	return ns
}

func (ss *StrSet) Difference(s *StrSet) *StrSet {
	ss.rLock()
	defer ss.rUnlock()

	var ns = NewStrSet(ss.block)
	for k := range ss.m {
		if !s.Contains(k) {
			ns.Add(k)
		}
	}
	return ns
}

func (ss *StrSet) ToStrList() []string {
	ss.rLock()
	defer ss.rUnlock()

	var stringList = make([]string, 0, ss.len())
	for k := range ss.m {
		stringList = append(stringList, k)
	}
	return stringList
}
