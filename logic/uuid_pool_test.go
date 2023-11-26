package logic

import (
	"math/rand"
	"testing"
	"time"

	"github.com/luci/go-render/render"
)

func TestUUIDPool(t *testing.T) {
	InitPool()
	time.Sleep(3 * time.Second)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				countRand := rand.Int63n(10000)
				bounds, err := MysqlUuidPool.GetUUIDBounds(1, countRand)
				if err != nil {
					t.Logf("num:%v GetUUIDBound err:%v \n", countRand, err)
				}
				t.Logf("num:%v GetUUIDBound:%v \n", countRand, render.Render(bounds))
				time.Sleep(300 * time.Millisecond)
			}
		}()
	}
	time.Sleep(100 * time.Second)
}
