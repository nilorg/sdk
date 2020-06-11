package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Uploader 上传
type Uploader interface {
	Upload(ctx context.Context, read io.Reader, parameters ...interface{}) (fullName string, err error)
}

// Storager 存储
type Storager interface {
	Uploader
	Download(ctx context.Context, write io.Writer, parameters ...interface{}) (results interface{}, err error)
	Remove(ctx context.Context, fullName string, parameters ...interface{}) (err error)
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

func (ds *DefaultStorage) filename(parameters ...interface{}) (filename string, err error) {
	if len(parameters) == 0 {
		err = errors.New("Please enter filename")
		return
	}
	switch parameters[0].(type) {
	case string:
		filename = parameters[0].(string)
	default:
		err = errors.New("filename parameter type error")
	}
	return
}

// Upload 上传
func (ds *DefaultStorage) Upload(_ context.Context, read io.Reader, parameters ...interface{}) (fullName string, err error) {
	var (
		filename string
	)
	filename, err = ds.filename(parameters...)
	if err != nil {
		return
	}
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
func (ds *DefaultStorage) Download(_ context.Context, dist io.Writer, parameters ...interface{}) (results interface{}, err error) {
	var (
		filename string
	)
	filename, err = ds.filename(parameters...)
	if err != nil {
		return
	}
	fullName := filepath.Join(ds.BasePath, filename)
	file, err := os.Open(fullName)
	if err != nil {
	}
	results, err = io.Copy(dist, file)
	return
}

// Remove 删除
func (ds *DefaultStorage) Remove(_ context.Context, fullName string, parameters ...interface{}) (err error) {
	err = os.Remove(fullName)
	return
}
