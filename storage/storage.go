package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// Uploader 上传
type Uploader interface {
	Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error)
}

// Downloader 下载
type Downloader interface {
	Download(ctx context.Context, write io.Writer, fullName string) (results interface{}, err error)
}

// Remover 删除
type Remover interface {
	Remove(ctx context.Context, filename string) (err error)
}

// Storager 存储
type Storager interface {
	Uploader
	Downloader
	Remover
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
func (ds *DefaultStorage) Upload(_ context.Context, read io.Reader, filename string) (fullName string, err error) {
	fullName = filepath.Join(ds.BasePath, filename)
	dir := filepath.Dir(fullName)
	_, dirErr := os.Stat(dir)
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			os.MkdirAll(dir, os.ModePerm)
		} else {
			err = dirErr
			return
		}
	}
	dst, err := os.Create(fullName)
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
func (ds *DefaultStorage) Download(_ context.Context, dist io.Writer, filename string) (results interface{}, err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	file, err := os.Open(fullName)
	if err != nil {
	}
	results, err = io.Copy(dist, file)
	return
}

// Remove 删除
func (ds *DefaultStorage) Remove(_ context.Context, filename string) (err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	err = os.Remove(fullName)
	return
}
