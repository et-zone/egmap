package egmap

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Map struct {
	mu sync.Mutex
	m  map[string]*entry
}

// An entry is a slot in the map corresponding to a particular key.
type entry struct {
	ep int64          //expire=0,never expire.
	p  unsafe.Pointer // *interface{}
}

func newEntry(i interface{}) *entry {
	return &entry{p: unsafe.Pointer(&i)}
}

func NewEMap() *Map {
	return &Map{m: map[string]*entry{}}
}

// Load returns the value stored in the map for a key, or nil if no
// value is present
func (m *Map) Load(key string) (value interface{}, ok bool) {
	m.mu.Lock()
	if e, ok := m.m[key]; ok {
		if e.ep > 0 && time.Now().Unix() >= e.ep {
			m.Delete(key)
			m.mu.Unlock()
			return nil, true
		}
		m.mu.Unlock()
		return e.load()
	}
	m.mu.Unlock()
	return nil, false
}

func (e *entry) load() (value interface{}, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil {
		return nil, false
	}
	return *(*interface{})(p), true
}

// Store sets the value for a key.
func (m *Map) Store(key string, value interface{}) {
	m.mu.Lock()
	if key != "" && value != nil {
		if e, ok := m.m[key]; ok && e != nil {
			e.tryStore(&value)
		} else {
			m.m[key] = newEntry(value)
		}
	}
	m.mu.Unlock()
}

// expire =0 ,expire>0
func (m *Map) StoreWithExpire(key string, value interface{}, expire int64) {
	m.mu.Lock()
	if key != "" && value != nil {
		if expire < 0 {
			panic("painc:Expire time must >= 0 ")
		}
		if e, ok := m.m[key]; ok && e != nil {
			e.tryStore(&value)
			if expire > 0 {
				e.ep = time.Now().Add(time.Second * time.Duration(expire)).Unix()
			}

		} else {
			m.m[key] = newEntry(value)
		}
	}
	m.mu.Unlock()
}

func (e *entry) tryStore(i *interface{}) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

// Delete deletes the value for a key.
func (m *Map) Delete(key string) {
	m.mu.Lock()
	if e, ok := m.m[key]; key != "" && ok && e != nil {
		e.delete()
		delete(m.m, key)
	}
	m.mu.Unlock()
}

func (e *entry) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

//资源清理，可选功能
func GC(m *Map) {
	go func() {
		for k, e := range m.m {
			if e.ep <= time.Now().Unix() {
				m.Delete(k)
			}
		}
	}()
}

func (m *Map) Range(f func(key string, value interface{}) bool) {
	for k, e := range m.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

//清空map -- 效率没有GC高，弃用 --
func (m *Map) CleanAll() {
	for k, _ := range m.m {
		m.Delete(k)
	}
}
