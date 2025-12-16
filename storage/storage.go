package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nilorg/sdk/mime"
)

// Uploader 上传
type Uploader interface {
	Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error)
}

// Downloader 下载
type Downloader interface {
	Download(ctx context.Context, write io.Writer, filename string) (info DownloadFileInfoer, err error)
}

// Remover 删除
type Remover interface {
	Remove(ctx context.Context, filename string) (err error)
}

type Exister interface {
	Exist(ctx context.Context, filename string) (exist bool, err error)
}

// Storager 存储
type Storager interface {
	Uploader
	Downloader
	Remover
	Exister
}

var _ Storager = (*DefaultStorage)(nil)

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
func (ds *DefaultStorage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
	if rename, ok := FromRenameContext(ctx); ok {
		filename = rename(filename)
	}
	saveFileName := filepath.Join(ds.BasePath, filename)
	dir := filepath.Dir(saveFileName)
	_, dirErr := os.Stat(dir)
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return
			}
		} else {
			err = dirErr
			return
		}
	}
	var dist *os.File
	dist, err = os.Create(saveFileName)
	if err != nil {
		return
	}
	defer dist.Close()
	_, err = io.Copy(dist, read)
	if err != nil {
		os.Remove(saveFileName)
		return
	}
	fullName = filename
	return
}

// Download 下载
func (ds *DefaultStorage) Download(ctx context.Context, dist io.Writer, filename string) (info DownloadFileInfoer, err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	file, err := os.Open(fullName)
	if err != nil {
		return
	}
	defer file.Close()

	md := Metadata{}
	if mimeType, exist := mime.Lookup(filepath.Ext(filename)); exist {
		md.Set("Content-Type", mimeType)
	}
	var (
		downloadFilename      string
		downloadFilenameExist bool
	)
	if downloadFilename, downloadFilenameExist = FromDownloadFilenameContext(ctx); !downloadFilenameExist {
		downloadFilename = filepath.Base(filename)
	}
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return
	}

	fileSize := fileInfo.Size()
	actualSize := fileSize

	// 检查是否有范围请求
	if rangeReq, rangeOk := FromRangeRequestContext(ctx); rangeOk {
		// 处理范围请求
		var startPos int64
		if rangeReq.Start < 0 {
			// Suffix Range: 读取最后 N 个字节
			suffixLen := -rangeReq.Start
			if suffixLen > fileSize {
				suffixLen = fileSize
			}
			startPos = fileSize - suffixLen
			actualSize = suffixLen
		} else {
			// Normal Range
			if (fileSize > 0 && rangeReq.Start >= fileSize) || (fileSize == 0 && rangeReq.Start > 0) {
				// 起始位置超出文件大小，返回错误
				err = fmt.Errorf("range start position %d exceeds file size %d", rangeReq.Start, fileSize)
				return
			}
			startPos = rangeReq.Start

			// 计算实际要读取的大小
			if rangeReq.IsToEnd() {
				// 读取到文件末尾
				actualSize = fileSize - startPos
			} else {
				// 读取指定范围
				if rangeReq.End >= fileSize {
					// 结束位置超出文件大小，调整到文件末尾
					actualSize = fileSize - startPos
				} else {
					actualSize = rangeReq.End - startPos + 1
				}
			}
		}

		// 设置文件读取位置到起始字节
		_, err = file.Seek(startPos, io.SeekStart)
		if err != nil {
			return
		}

		// 使用LimitReader限制读取的字节数
		limitedReader := io.LimitReader(file, actualSize)
		_, err = io.Copy(dist, limitedReader)
	} else {
		// 没有范围请求，读取整个文件
		_, err = io.Copy(dist, file)
	}

	info = &downloadFileInfo{
		size:     actualSize, // 返回实际传输的大小
		metadata: md,
		filename: downloadFilename,
	}
	if downloadBefore, downloadBeforeExist := FromDownloadBeforeContext(ctx); downloadBeforeExist {
		downloadBefore(info)
	}

	if err != nil {
		info = nil
	}
	return
}

// Remove 删除
func (ds *DefaultStorage) Remove(_ context.Context, filename string) (err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	err = os.Remove(fullName)
	return
}

// Exist 是否存在
func (ds *DefaultStorage) Exist(_ context.Context, filename string) (exist bool, err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	_, err = os.Stat(fullName)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	exist = true
	return
}

// DownloadFileInfoer 下载file信息接口
type DownloadFileInfoer interface {
	// Size 文件大小
	Size() int64
	// Filename 文件名
	Filename() string
	// Metadata 获取元数据
	Metadata() Metadata
}

type downloadFileInfo struct {
	size     int64
	filename string
	metadata Metadata
}

func (dfi *downloadFileInfo) Size() int64 {
	return dfi.size
}

func (dfi *downloadFileInfo) Filename() string {
	return dfi.filename
}

func (dfi *downloadFileInfo) Metadata() Metadata {
	return dfi.metadata
}
