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
