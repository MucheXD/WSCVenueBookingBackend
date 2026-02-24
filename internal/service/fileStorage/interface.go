package fileStorage

import (
	"context"
	"io"
)

type FileStorageSvc interface {
	// Upload 上传文件流
	Upload(ctx context.Context, fileToken string, fileReader io.Reader) error

	// Download 下载文件流
	Download(ctx context.Context, filetoken string) (io.ReadCloser, error)

	// Delete 删除文件
	Delete(ctx context.Context, filetoken string) error
}
