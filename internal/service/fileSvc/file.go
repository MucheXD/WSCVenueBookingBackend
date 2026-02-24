package fileSvc

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/fileStorage"
	"gorm.io/gorm"
)

const (
	StorageTypeLocal = 1

	tokenTimePartLength = 8
	fileTokenLength     = 60
)

var (
	sha256HexRegexp  = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
	fileTokenRegexp  = regexp.MustCompile(`^[A-Z2-7]{60}$`)
	base32HexEncoder = base32.StdEncoding.WithPadding(base32.NoPadding)

	storageInit sync.Once
	localStore  fileStorage.FileStorageSvc
)

func UploadFile(ctx context.Context, size int64, claimedHash string, fileReader io.Reader) (string, error) {
	if fileReader == nil {
		return "", ErrFileParamInvalid
	}
	if size < 0 {
		return "", ErrFileSizeInvalid
	}

	// 将前端传入哈希统一为大写并去除空白、验证格式
	nClaimedHash := strings.ToUpper(strings.TrimSpace(claimedHash))
	if !sha256HexRegexp.MatchString(nClaimedHash) {
		return "", ErrFileHashInvalid
	}

	// 秒传：查询是否已存在相同哈希的文件，避免重复上传
	existing, err := repository.GetFileObjectByHash(nClaimedHash)
	if err == nil {
		if existing.FileSize != size {
			return "", ErrFileHashSizeConflict
		}
		return existing.FileToken, nil // 已存在相同文件，直接返回 token
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("%w: %w", ErrFileQueryInDB, err)
	}

	fileToken, err := buildFileToken(time.Now().Unix(), nClaimedHash)
	if err != nil {
		return "", err
	}

	storageSvc, err := resolveStorageByType(StorageTypeLocal)
	if err != nil {
		return "", err
	}

	verifyReader := &hashVerifyReader{
		source: fileReader,
		hasher: sha256.New(),
	}

	// 开始上传行为
	if err := storageSvc.Upload(ctx, fileToken, verifyReader); err != nil {
		return "", fmt.Errorf("%w: %w", ErrFileUploadToStorage, err)
	}

	// 验证文件哈希与大小是否与前端声明一致
	actualSize := verifyReader.size
	if actualSize != size {
		_ = storageSvc.Delete(ctx, fileToken)
		return "", fmt.Errorf("%w: expect=%d, actual=%d", ErrFileSizeMismatch, size, actualSize)
	}

	actualHash := strings.ToUpper(hex.EncodeToString(verifyReader.hasher.Sum(nil)))
	if actualHash != nClaimedHash {
		_ = storageSvc.Delete(ctx, fileToken)
		return "", ErrFileHashMismatch
	}

	// 创建数据库记录
	modelFile := &models.FileObject{
		FileToken:   fileToken,
		FileHash:    nClaimedHash,
		FileSize:    size,
		StorageType: StorageTypeLocal,
	}
	if err := repository.CreateFileObject(modelFile); err != nil {
		_ = storageSvc.Delete(ctx, fileToken)

		concurrent, qErr := repository.GetFileObjectByHash(nClaimedHash)
		if qErr == nil {
			if concurrent.FileSize != size {
				return "", ErrFileHashSizeConflict
			}
			return concurrent.FileToken, nil
		}
		return "", fmt.Errorf("%w: %w", ErrFileCreateInDB, err)
	}

	return fileToken, nil
}

func DownloadFile(ctx context.Context, fileToken string) (io.ReadCloser, int64, error) {
	if !isFileTokenValid(fileToken) {
		return nil, 0, ErrFileTokenInvalid
	}

	modelFile, err := repository.GetFileObjectByToken(fileToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrFileNotFound
		}
		return nil, 0, fmt.Errorf("%w: %w", ErrFileQueryInDB, err)
	}

	storageSvc, err := resolveStorageByType(modelFile.StorageType)
	if err != nil {
		return nil, 0, err
	}

	reader, err := storageSvc.Download(ctx, fileToken)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, 0, ErrFileNotFound
		}
		return nil, 0, fmt.Errorf("%w: %w", ErrFileDownloadFromStorage, err)
	}
	return reader, modelFile.FileSize, nil
}

// 解析文件存储类型，目前仅支持本地存储
func resolveStorageByType(storageType int) (fileStorage.FileStorageSvc, error) {
	if storageType != StorageTypeLocal {
		return nil, fmt.Errorf("%w: %d", ErrFileStorageTypeUnsupported, storageType)
	}
	return getLocalStorage()
}

func getLocalStorage() (fileStorage.FileStorageSvc, error) {
	storageInit.Do(func() {
		rootPath := strings.TrimSpace(config.Config.GetString("storage.local.root_path"))
		if rootPath == "" {
			rootPath = "./storage"
		}
		localStore = fileStorage.NewLocalStorage(rootPath)
	})
	if localStore == nil {
		return nil, ErrFileStorageNotConfigured
	}
	return localStore, nil
}

// buildFileToken 生成文件token，格式为：{8 chars time part}{52 chars hash part}
func buildFileToken(unixSeconds int64, sha256HexUpper string) (string, error) {
	hashBytes, err := hex.DecodeString(sha256HexUpper)
	if err != nil {
		return "", ErrFileHashInvalid
	}
	timePart := encodeUnixToBase32Fixed(unixSeconds, tokenTimePartLength)
	hashPart := base32HexEncoder.EncodeToString(hashBytes)
	fileToken := timePart + hashPart
	if len(fileToken) != fileTokenLength {
		return "", ErrFileTokenInvalid
	}
	return fileToken, nil
}

// encodeUnixToBase32Fixed 将unix时间戳编码为固定长度的base32字符串
func encodeUnixToBase32Fixed(unixSeconds int64, width int) string {
	if unixSeconds < 0 {
		unixSeconds = 0
	}
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	buffer := make([]byte, width)
	value := unixSeconds
	for i := width - 1; i >= 0; i-- {
		buffer[i] = charset[value%32]
		value /= 32
	}
	return string(buffer)
}

func isFileTokenValid(fileToken string) bool {
	return fileTokenRegexp.MatchString(fileToken)
}

// hashVerifyReader 包装一个 io.Reader，在读取数据的同时计算哈希和记录大小
// 使用时，按照 io.Reader 方式读取，完成后可通过 hasher.Sum(nil) 获取哈希结果，通过 size 字段获取总大小
type hashVerifyReader struct {
	source io.Reader
	hasher hash.Hash
	size   int64
}

func (r *hashVerifyReader) Read(p []byte) (int, error) {
	n, err := r.source.Read(p)
	if n > 0 {
		r.size += int64(n)
		_, _ = r.hasher.Write(p[:n])
	}
	return n, err
}
