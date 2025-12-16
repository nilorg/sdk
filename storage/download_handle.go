package storage

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nilorg/sdk/mime"
)

// DownloadStorager 下载接口
type DownloadStorager interface {
	Downloader
	DispositionType() string
}

// DefaultDownloadStorage 默认下载
type DefaultDownloadStorage struct {
	*DefaultStorage
	dispositionType string
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
func DownloadHandle(ctx context.Context, resp http.ResponseWriter, ds DownloadStorager, filename string) (err error) {
	if ds == nil {
		ds = NewDefaultDownloadStorage("attachment")
	}

	if _, ok := FromDownloadBeforeContext(ctx); !ok {
		ctx = NewDownloadBeforeContext(ctx, func(info DownloadFileInfoer) {
			md := info.Metadata()
			wr := false
			if md != nil {
				contextType := md.Get("Content-Type")
				if contextType != "" {
					resp.Header().Set("Content-Type", contextType)
					wr = true
				}
			}
			if !wr {
				var contentType string
				contentType, err = mime.DetectContentType(info.Filename())
				if err != nil {
					return
				}
				resp.Header().Set("Content-Type", contentType)
			}
			dispositionType := ds.DispositionType()
			if dispositionType == "" || (dispositionType != "inline" && dispositionType != "attachment") {
				dispositionType = "attachment"
			}
			resp.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", dispositionType, info.Filename()))
		})
	}
	_, err = ds.Download(ctx, resp, filename)
	if err != nil {
		return
	}
	resp.WriteHeader(http.StatusOK)
	return
}

// DownloadHandleWithRange 支持范围请求的下载处理
func DownloadHandleWithRange(ctx context.Context, req *http.Request, resp http.ResponseWriter, ds DownloadStorager, filename string) (err error) {
	if ds == nil {
		ds = NewDefaultDownloadStorage("attachment")
	}

	// 解析 Range 头
	rangeHeader := req.Header.Get("Range")
	if rangeHeader != "" {
		start, end, parseErr := parseRangeHeader(rangeHeader)
		if parseErr == nil {
			ctx = NewRangeRequestContext(ctx, start, end)
		}
	}

	// 检查是否有范围请求
	rangeReq, hasRange := FromRangeRequestContext(ctx)
	isPartialContent := hasRange && rangeReq != nil

	if _, ok := FromDownloadBeforeContext(ctx); !ok {
		ctx = NewDownloadBeforeContext(ctx, func(info DownloadFileInfoer) {
			md := info.Metadata()
			wr := false
			if md != nil {
				contextType := md.Get("Content-Type")
				if contextType != "" {
					resp.Header().Set("Content-Type", contextType)
					wr = true
				}
			}
			if !wr {
				var contentType string
				contentType, err = mime.DetectContentType(info.Filename())
				if err != nil {
					return
				}
				resp.Header().Set("Content-Type", contentType)
			}
			dispositionType := ds.DispositionType()
			if dispositionType == "" || (dispositionType != "inline" && dispositionType != "attachment") {
				dispositionType = "attachment"
			}
			resp.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", dispositionType, info.Filename()))

			// 设置Accept-Ranges头，表示支持范围请求
			resp.Header().Set("Accept-Ranges", "bytes")

			if isPartialContent {
				// 设置Content-Range头
				if rangeReq.IsToEnd() {
					resp.Header().Set("Content-Range", fmt.Sprintf("bytes %d-*/*", rangeReq.Start))
				} else {
					resp.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/*", rangeReq.Start, rangeReq.End))
				}
			}
		})
	}

	// 在写入响应体之前设置状态码
	if isPartialContent {
		resp.WriteHeader(http.StatusPartialContent)
	} else {
		resp.WriteHeader(http.StatusOK)
	}

	_, err = ds.Download(ctx, resp, filename)
	if err != nil {
		return
	}

	return
}

// parseRangeHeader 解析 HTTP Range 头
// 支持格式: "bytes=start-end" 或 "bytes=start-"
func parseRangeHeader(rangeHeader string) (start, end int64, err error) {
	// Range: bytes=0-1023
	// Range: bytes=1024-
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		err = fmt.Errorf("unsupported range format: %s", rangeHeader)
		return
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")

	// 支持简单的 start-end 格式
	parts := strings.SplitN(rangeSpec, "-", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid range format: %s", rangeSpec)
		return
	}

	// 解析起始位置
	if parts[0] != "" {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("invalid range start: %s", parts[0])
			return
		}
	} else {
		// 如果起始位置为空，必须有结束位置，表示 suffix-byte-range-spec
		// bytes=-500 -> start=-500 (internal representation)
		if parts[1] == "" {
			err = fmt.Errorf("invalid range format: %s", rangeSpec)
			return
		}
		var suffixLen int64
		suffixLen, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("invalid range suffix length: %s", parts[1])
			return
		}
		start = -suffixLen
		end = -1 // Suffix range implies to end
		return
	}

	// 解析结束位置
	if parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("invalid range end: %s", parts[1])
			return
		}
	} else {
		// 没有指定结束位置，表示到文件末尾
		end = -1
	}

	return
}
