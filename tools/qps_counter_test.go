package tools

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestQPSCounter(t *testing.T) {
	counter := NewQPSCounter(30)
	for i := 1; i < 3; i++ {
		go func() {
			for {
				counter.AddCount(50)
				time.Sleep(1000 * time.Millisecond)
			}
		}()
	}
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 1,count:", counter.SumCount())
		}
	}()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 2,count:", counter.SumCount())
		}
	}()

	time.Sleep(1000 * time.Second)
}

func TestQPSCounterHigh(t *testing.T) {
	counter := NewQPSCounter(30)
	for i := 1; i < 30; i++ {
		go func() {
			for {
				counter.AddCount(50)
				time.Sleep(1000 * time.Millisecond)
			}
		}()
	}
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 1,count:", counter.SumCount())
		}
	}()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 2,count:", counter.SumCount())
		}
	}()

	time.Sleep(1000 * time.Second)
}

func TestQPSCounterSand(t *testing.T) {
	counter := NewQPSCounter(30)
	for i := 1; i < 30; i++ {
		go func() {
			for {
				counter.AddCount(rand.Int63n(1000))
				time.Sleep(1000 * time.Millisecond)
			}
		}()
	}
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 1,count:", counter.SumCount())
		}
	}()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("first 2,count:", counter.SumCount())
		}
	}()

	time.Sleep(1000 * time.Second)
}
