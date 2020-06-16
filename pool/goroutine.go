package pool

import (
	"sync"
)

// Goroutine 池
type Goroutine struct {
	queue chan int
	wg    *sync.WaitGroup
}

// NewGoroutine 创建池大小
func NewGoroutine(size int) *Goroutine {
	if size <= 0 {
		size = 1
	}
	return &Goroutine{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add 添加
func (p *Goroutine) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done 完成
func (p *Goroutine) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait 等待
func (p *Goroutine) Wait() {
	p.wg.Wait()
}
