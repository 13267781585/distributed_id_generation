package logic

import (
	"math/rand"
	"testing"
	"time"

	"github.com/luci/go-render/render"
)

func TestBufferSingleGetHighQPS(t *testing.T) {
	buffer := NewUUIDBuffer(NewDefaultMysqlConfig())
	buffer.Start()
	time.Sleep(5 * time.Second)

	GetUUIDS(&Param{
		Test:        t,
		Buffer:      buffer,
		Count:       100000,
		Interval:    time.Second,
		ConsumerNum: 1,
	})

	time.Sleep(200 * time.Second)
}

func TestBufferSingleGetLowQPS(t *testing.T) {
	buffer := NewUUIDBuffer(NewDefaultMysqlConfig())
	buffer.Start()
	time.Sleep(5 * time.Second)

	GetUUIDS(&Param{
		Test:        t,
		Buffer:      buffer,
		Count:       1000,
		Interval:    time.Second,
		ConsumerNum: 1,
	})
	time.Sleep(200 * time.Second)
}

func TestBufferMultiGetLowQPS(t *testing.T) {
	buffer := NewUUIDBuffer(NewDefaultMysqlConfig())
	buffer.Start()
	time.Sleep(5 * time.Second)

	GetUUIDS(&Param{
		Test:        t,
		Buffer:      buffer,
		Count:       1000,
		Interval:    time.Second,
		ConsumerNum: 10,
	})
	time.Sleep(200 * time.Second)
}

func TestBufferMultiGetHighQPS(t *testing.T) {
	buffer := NewUUIDBuffer(NewDefaultMysqlConfig())
	buffer.Start()
	time.Sleep(5 * time.Second)

	GetUUIDS(&Param{
		Test:        t,
		Buffer:      buffer,
		Count:       10000,
		Interval:    time.Second,
		ConsumerNum: 10,
	})
	time.Sleep(200 * time.Second)
}

type Param struct {
	Test        *testing.T
	Buffer      *UUIDBuffer
	Count       int64
	Interval    time.Duration
	ConsumerNum int64
}

func GetUUIDS(param *Param) {
	for i := int64(0); i < param.ConsumerNum; i++ {
		go func(num int64) {
			for {
				countRand := rand.Int63n(param.Count)
				bounds, err := param.Buffer.GetUUIDBound(countRand)
				if err != nil {
					param.Test.Logf("num:%v GetUUIDBound err:%v \n", num, err)
				}
				param.Test.Logf("num:%v GetUUIDBound:%v \n", num, render.Render(bounds))
				time.Sleep(param.Interval)
			}
		}(i)
	}

}
