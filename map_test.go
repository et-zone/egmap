package egmap

import (
	"fmt"

	"testing"

	emap "github.com/et-zone/egmap"
)

type V struct {
	Data string
}

func Test_emap(t *testing.T) {
	mp := emap.NewEMap()
	mp.Store("aa", &V{"aa"})
	mp.Store("bb", &V{"bb"})
	mp.Store("cc", &V{"cc"})
	v, ok := mp.Load("aa")
	if ok {
		fmt.Println(v.(*V).Data)
	}
	mp.Store("bb", &V{"ff"})
	mp.Range(func(k string, v interface{}) bool {
		fmt.Println(v)
		// mp.Delete(k) //del
		return true
	})

	go emap.GC(mp)

}
