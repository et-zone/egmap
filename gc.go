package egmap

// import (
// 	"sync/atomic"
// 	"unsafe"
// )

//gc 在随机range map时，很大概率会命中再set的key，导致读写同一个key(加锁无效)，因此，弃用gc
// 系统自带的sync.map  ,有一个脏字典，可以读写分离，不会导致这个问题
// var regist *[]unsafe.Pointer
// var gcFlag = false

// //资源gc，可选功能
// func GCRun() {
// 	if gcFlag {
// 		return
// 	}
// 	go func() {
// 		for _, w := range *regist {
// 			if w != nil {
// 				m := *(*Map)(w)
// 				m.delsWithExpire()
// 			}
// 		}

// 	}()
// 	gcFlag = true
// }

// //注册需要gc的map
// func RegistGC(m *Map) {
// 	if regist == nil {
// 		r := make([]unsafe.Pointer, 0)
// 		regist = &r
// 	}
// 	tmpp := unsafe.Pointer(&m)
// 	for _, p := range *regist {

// 		if atomic.LoadPointer(&p) == tmpp {
// 			return
// 		}
// 	}
// 	*regist = append(*regist, tmpp)
// }
