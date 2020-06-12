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
