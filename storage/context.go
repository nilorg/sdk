package storage

import (
	"context"
)

type downloadFilenameKey struct{}

// NewDownloadFilenameContext ...
func NewDownloadFilenameContext(ctx context.Context, downloadFilename string) context.Context {
	return context.WithValue(ctx, downloadFilenameKey{}, downloadFilename)
}

// FromDownloadFilenameContext ...
func FromDownloadFilenameContext(ctx context.Context) (downloadFilename string, ok bool) {
	downloadFilename, ok = ctx.Value(downloadFilenameKey{}).(string)
	return
}

type downloadBeforeKey struct{}

// DownloadBefore 下载之前
type DownloadBefore func(info DownloadFileInfoer)

// NewDownloadBeforeContext ...
func NewDownloadBeforeContext(ctx context.Context, f DownloadBefore) context.Context {
	return context.WithValue(ctx, downloadBeforeKey{}, f)
}

// FromDownloadBeforeContext ...
func FromDownloadBeforeContext(ctx context.Context) (f DownloadBefore, ok bool) {
	f, ok = ctx.Value(downloadBeforeKey{}).(DownloadBefore)
	return
}

type renameKey struct{}

// Rename 重命名
type Rename func(filename string) string

// NewRenameContext ...
func NewRenameContext(ctx context.Context, f Rename) context.Context {
	return context.WithValue(ctx, renameKey{}, f)
}

// FromRenameContext ...
func FromRenameContext(ctx context.Context) (f Rename, ok bool) {
	f, ok = ctx.Value(renameKey{}).(Rename)
	return
}

// ---

// RangeRequest 表示Range请求
type RangeRequest struct {
	Start int64 // 开始字节位置（包含），从0开始
	End   int64 // 结束字节位置（包含），<=0表示到文件末尾
}

// IsValid 检查Range请求是否有效
func (r *RangeRequest) IsValid() bool {
	if r.Start < 0 {
		return false
	}
	// End <= 0 表示到文件末尾，这是有效的
	// End > 0 且 End >= Start 也是有效的
	if r.End > 0 && r.End < r.Start {
		return false
	}
	return true
}

// IsToEnd 检查是否是到文件末尾的请求
func (r *RangeRequest) IsToEnd() bool {
	return r.End <= 0
}

// Length 计算请求的字节长度，如果是到文件末尾则返回-1
func (r *RangeRequest) Length() int64 {
	if r.End <= 0 {
		return -1
	}
	return r.End - r.Start + 1
}

type rangeRequestKey struct{}

// NewRangeRequestContext 创建包含Range请求的context
// start: 开始字节位置（>=0）
// end: 结束字节位置（>start，或<=0表示到文件末尾）
// 例如：start=0, end=99 表示读取前100个字节（0到99）
//
//	start=100, end=0 表示从第101个字节读取到文件末尾
//	start=0, end=-1 表示读取整个文件
func NewRangeRequestContext(ctx context.Context, start, end int64) context.Context {
	rangeReq := &RangeRequest{
		Start: start,
		End:   end,
	}

	// 验证Range请求的有效性
	if !rangeReq.IsValid() {
		// 对于无效的Range请求，可以选择panic或返回原context
		// 这里选择返回原context，调用者需要检查FromRangeRequestContext的返回值
		return ctx
	}

	return context.WithValue(ctx, rangeRequestKey{}, rangeReq)
}

// FromRangeRequestContext 从context中获取Range请求
func FromRangeRequestContext(ctx context.Context) (rangeReq *RangeRequest, ok bool) {
	rangeReq, ok = ctx.Value(rangeRequestKey{}).(*RangeRequest)
	return
}
