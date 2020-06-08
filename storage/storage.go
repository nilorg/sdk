package storage

import (
	"io"
	"os"
	"path/filepath"
)

// Storager 存储
type Storager interface {
	Upload(read io.Reader, filename string) (fullPath string, err error)
	Download(write io.Writer, filename string) (err error)
}

// DefaultStorage 默认存储
type DefaultStorage struct {
	BasePath string
}

// NewDefaultStorage 创建默认存储
func NewDefaultStorage() *DefaultStorage {
	return &DefaultStorage{
		BasePath: "./",
	}
}

// Upload 上传
func (ds *DefaultStorage) Upload(read io.Reader, filename string) (fullPath string, err error) {
	fullPath = filepath.Join(ds.BasePath, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return
	}
	_, err = io.Copy(dst, read)
	if err != nil {
		return
	}
	return
}

// Download 下载
func (ds *DefaultStorage) Download(dist io.Writer, filename string) (err error) {
	fullPath := filepath.Join(ds.BasePath, filename)
	file, err := os.Open(fullPath)
	if err != nil {
	}
	_, err = io.Copy(dist, file)
	return
}
