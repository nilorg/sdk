package pool

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)

type factory func() (io.Closer, error)

// Pooler 连接池接口
type Pooler interface {
	Get() (io.Closer, error) // 获取资源
	Put(io.Closer) error     // 释放资源
	Close(io.Closer) error   // 关闭资源
	Shutdown() error         // 关闭池
}

// CommonPool 通用连接池
type CommonPool struct {
	sync.Mutex
	pool    chan io.Closer
	maxOpen int     // 池中最大资源数
	numOpen int     // 当前池中资源数
	minOpen int     // 池中最少资源数
	closed  bool    // 池是否已关闭
	factory factory // 创建连接的方法
}

// NewCommonPool 创建通用连接池
func NewCommonPool(minOpen, maxOpen int, factory factory) (*CommonPool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &CommonPool{
		maxOpen: maxOpen,
		minOpen: minOpen,
		factory: factory,
		pool:    make(chan io.Closer, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		closer, err := factory()
		if err != nil {
			continue
		}
		p.numOpen++
		p.pool <- closer
	}
	return p, nil
}

// Get 获取资源
func (p *CommonPool) Get() (io.Closer, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	for {
		closer, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		return closer, nil
	}
}

func (p *CommonPool) getOrCreate() (io.Closer, error) {
	select {
	case closer := <-p.pool:
		return closer, nil
	default:
	}
	p.Lock()
	defer p.Unlock()
	if p.numOpen >= p.maxOpen {
		closer := <-p.pool
		return closer, nil
	}
	// 新建连接
	closer, err := p.factory()
	if err != nil {
		return nil, err
	}
	p.numOpen++
	return closer, nil
}

// Put 释放单个资源到连接池
func (p *CommonPool) Put(closer io.Closer) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	defer p.Unlock()
	p.pool <- closer
	return nil
}

// Close 关闭单个资源
func (p *CommonPool) Close(closer io.Closer) error {
	p.Lock()
	defer p.Unlock()
	closer.Close()
	p.numOpen--
	return nil
}

// Shutdown 关闭连接池，释放所有资源
func (p *CommonPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	defer p.Unlock()
	close(p.pool)
	for closer := range p.pool {
		closer.Close()
		p.numOpen--
	}
	p.closed = true
	return nil
}
