package fileStorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地文件存储 FileStorage 实现
type LocalStorage struct {
	RootPath string
}

func NewLocalStorage(rootPath string) *LocalStorage {
	return &LocalStorage{
		RootPath: rootPath,
	}
}

func (l *LocalStorage) Upload(ctx context.Context, fileToken string, fileReader io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if fileToken == "" {
		return fmt.Errorf("file token is empty")
	}

	targetPath := l.filePath(fileToken)
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	tempPath := targetPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(tempPath)
	}()

	if _, err := io.Copy(file, fileReader); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, targetPath); err != nil {
		return err
	}
	return nil
}

func (l *LocalStorage) Download(ctx context.Context, filetoken string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return os.Open(l.filePath(filetoken))
}

func (l *LocalStorage) Delete(ctx context.Context, filetoken string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := os.Remove(l.filePath(filetoken))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (l *LocalStorage) filePath(fileToken string) string {
	folder := fileToken
	if len(fileToken) >= 4 {
		folder = fileToken[:4]
	}
	return filepath.Join(l.RootPath, folder, fileToken)
}
