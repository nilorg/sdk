package pool

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// LifeCycleCommonPool 生命周期通用连接池
type LifeCycleCommonPool struct {
	sync.Mutex
	ticker       *time.Ticker
	tickerCtx    context.Context
	tickerCancel context.CancelFunc
	pool         []io.Closer
	poolIndex    int           // 当前池中资源数
	numOpen      int           // 当前池中资源数
	maxOpen      int           // 池中最大资源数
	minOpen      int           // 池中最少资源数
	closed       bool          // 池是否已关闭
	factory      factory       // 创建连接的方法
	maxLifeTime  time.Duration // 最大生命周期
}

// Connection 生命周期
type Connection struct {
	io.Closer
	time time.Time
	use  bool
}

// NewLifeCycleCommonPool 创建通用连接池
func NewLifeCycleCommonPool(minOpen, maxOpen int, maxLifeTime time.Duration, factory factory) (*LifeCycleCommonPool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}
	tickerCtx, tickerCancel := context.WithCancel(context.Background())
	p := &LifeCycleCommonPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifeTime: maxLifeTime,
		factory:     factory,
		ticker:      time.NewTicker(time.Second * 1),
		tickerCtx:   tickerCtx, tickerCancel: tickerCancel,
		pool:      make([]io.Closer, 0),
		poolIndex: 0,
	}
	p.inactive()

	for i := 0; i < minOpen; i++ {
		closer, err := factory()
		if err != nil {
			continue
		}
		p.numOpen++
		p.pool = append(p.pool, &Connection{
			time:   time.Now(),
			Closer: closer,
			use:    false,
		})
	}
	return p, nil
}

// Get 获取资源
func (p *LifeCycleCommonPool) Get() (io.Closer, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	for {
		lc, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		if lc == nil && err == nil {
			time.Sleep(time.Millisecond)
			continue
		}
		return lc, nil
	}
}

func (p *LifeCycleCommonPool) getOrCreate() (io.Closer, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.pool) > 0 {
		// if p.poolIndex >= len(p.pool) || p.poolIndex >= p.maxOpen {
		if p.poolIndex >= len(p.pool) {
			p.poolIndex = 0
		}
		fmt.Printf("len: %d, index: %d, numOpen: %d\n", len(p.pool), p.poolIndex, p.numOpen)

		lc := p.pool[p.poolIndex].(*Connection)
		lc.use = true
		lc.time = time.Now()

		p.pool = append(p.pool[:p.poolIndex], p.pool[p.poolIndex+1:]...)
		p.poolIndex++

		return lc, nil
	}

	if p.numOpen < p.maxOpen {
		// 新建连接
		closer, err := p.factory()
		if err != nil {
			return nil, err
		}
		p.numOpen++
		return &Connection{
			Closer: closer,
			time:   time.Now(),
			use:    true,
		}, nil
	}

	return nil, nil
}

// Put 释放单个资源到连接池
func (p *LifeCycleCommonPool) Put(closer io.Closer) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	defer p.Unlock()
	if len(p.pool) == p.maxOpen {
		return ErrInvalidConfig
	}
	lc := closer.(*Connection)
	lc.use = false
	lc.time = time.Now()
	p.pool = append(p.pool, lc)
	return nil
}

// Close 关闭单个资源
func (p *LifeCycleCommonPool) Close(closer io.Closer) error {
	p.Lock()
	defer p.Unlock()
	err := closer.Close()
	p.numOpen--
	return err
}

// Shutdown 关闭连接池，释放所有资源
func (p *LifeCycleCommonPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	defer p.Unlock()
	// 停止自动检查不活跃的资源
	p.ticker.Stop()
	for i := 0; i < len(p.pool); i++ {
		lc := p.pool[i].(*Connection)
		lc.use = false
		err := lc.Closer.Close()
		if err != nil {
			return err
		}
		p.numOpen--
	}
	p.tickerCancel()
	p.closed = true
	p.pool = make([]io.Closer, 0)
	p.poolIndex = 0
	return nil
}

// inactive 释放不活跃资源
func (p *LifeCycleCommonPool) inactive() {
	go func() {
		for {
			select {
			case <-p.ticker.C:
				p.Lock()
				for i := 0; i < len(p.pool); i++ {
					lc := p.pool[i].(*Connection)
					d := time.Since(lc.time)
					if !lc.use && d >= p.maxLifeTime {
						lc.Closer.Close()
						p.pool = append(p.pool[:i], p.pool[i+1:]...)
						p.numOpen--
					}
				}
				p.Unlock()
			case <-p.tickerCtx.Done():
				return
			}
		}
	}()
}
