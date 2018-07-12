package goutil

import (
	"sync"
	"time"
)

// InstantTicker 类似于 time.Ticker。不同在于会立马发送当前时间。
type InstantTicker struct {
	C  <-chan time.Time
	tt *time.Ticker

	stopOnce   sync.Once // 用以防止多次 Stop 操作关闭 terminated 而 panic。
	terminated chan struct{}
}

// NewInstantTicker 创建一个新的 InstantTicker。立马向 InstantTicker 的 channel 中
// 发送当前的时间，并以 d 为时间间隔不断地向其中发送时间。
func NewInstantTicker(d time.Duration) *InstantTicker {
	c := make(chan time.Time, 1)
	it := &InstantTicker{
		C:          c,
		tt:         time.NewTicker(d),
		terminated: make(chan struct{}),
	}

	go func() {
		c <- time.Now()

	LOOP:
		for {
			select {
			case <-it.terminated:
				it.tt.Stop()
				break LOOP
			case t := <-it.tt.C:
				c <- t
			}
		}
	}()

	return it
}

// Stop 停止当前的 InstantTicker，但不会关闭其中的 channel 以防止错误行为。
func (it *InstantTicker) Stop() {
	it.stopOnce.Do(func() {
		close(it.terminated)
	})
}
