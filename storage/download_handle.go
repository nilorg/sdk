package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nilorg/sdk/mime"
)

// DownloadStorager 下载接口
type DownloadStorager interface {
	Downloader
	FormatWriteFilename(filename interface{}) string
	DispositionType() string
}

// DefaultDownloadStorage 默认下载
type DefaultDownloadStorage struct {
	*DefaultStorage
	dispositionType string
}

// FormatWriteFilename 格式化写入文件名
func (*DefaultDownloadStorage) FormatWriteFilename(filename interface{}) string {
	return filename.(string)
}

// DispositionType Disposition Type
func (d *DefaultDownloadStorage) DispositionType() string {
	return d.dispositionType
}

// NewDefaultDownloadStorage 创建默认存储
func NewDefaultDownloadStorage(dispositionType string) *DefaultDownloadStorage {
	return &DefaultDownloadStorage{
		DefaultStorage:  NewDefaultStorage(),
		dispositionType: dispositionType,
	}
}

// DownloadHandle 下载处理
func DownloadHandle(ctx context.Context, resp http.ResponseWriter, ds DownloadStorager, fullName string) (err error) {
	if ds == nil {
		ds = NewDefaultDownloadStorage("attachment")
	}
	var results interface{}
	results, err = ds.Download(ctx, resp, fullName)
	if err != nil {
		return
	}
	dispositionType := ds.DispositionType()
	if dispositionType == "" || (dispositionType != "inline" && dispositionType != "attachment") {
		dispositionType = "attachment"
	}
	writeFilename := ds.FormatWriteFilename(results)
	resp.Header().Add("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", dispositionType, writeFilename))

	if mimeType, exist := mime.Lookup(writeFilename); exist {
		resp.Header().Add("Content-Type", mimeType)
	}
	return
}
