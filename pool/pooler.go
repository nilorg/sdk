package pool

import (
	"errors"
	"io"
)

var (
	// ErrInvalidConfig invalid pool config
	ErrInvalidConfig = errors.New("invalid pool config")
	// ErrPoolClosed pool closed
	ErrPoolClosed = errors.New("pool closed")
)

type factory func() (io.Closer, error)

// Pooler 连接池接口
type Pooler interface {
	Get() (io.Closer, error) // 获取资源
	Put(io.Closer) error     // 释放资源
	Close(io.Closer) error   // 关闭资源
	Shutdown() error         // 关闭池
}
