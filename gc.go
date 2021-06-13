package egmap

import (
	"sync/atomic"
	"unsafe"
)

var regist *[]unsafe.Pointer
var gcFlag = false

//资源gc，可选功能
func GCRun() {
	if gcFlag {
		return
	}
	go func() {
		for _, w := range *regist {
			if w != nil {
				m := *(*Map)(w)
				m.delsWithExpire()
			}
		}

	}()
	gcFlag = true
}

//注册需要gc的map
func RegistGC(m *Map) {
	if regist == nil {
		r := make([]unsafe.Pointer, 0)
		regist = &r
	}
	tmpp := unsafe.Pointer(&m)
	for _, p := range *regist {

		if atomic.LoadPointer(&p) == tmpp {
			return
		}
	}
	*regist = append(*regist, tmpp)
}
