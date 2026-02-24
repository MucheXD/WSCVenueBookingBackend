package fileCtrl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/fileSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// UploadFileHandler 上传文件
// POST /api/file (multipart/form-data)
func UploadFileHandler(c *gin.Context) {
	sizeRaw := c.PostForm("size")
	hashRaw := c.PostForm("hash")
	if sizeRaw == "" || hashRaw == "" {
		apiException.AbortWithException(c, apiException.ParamError, errors.New("size and hash are required"))
		return
	}

	size, err := strconv.ParseInt(sizeRaw, 10, 64)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if fileHeader.Size >= 0 && fileHeader.Size != size {
		apiException.AbortWithException(c, apiException.ParamError, errors.New("file size mismatch with form size"))
		return
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	defer func() { _ = fileReader.Close() }()

	fileToken, err := fileSvc.UploadFile(c.Request.Context(), size, hashRaw, fileReader)
	if err != nil {
		handleFileServiceError(c, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]string{"filetoken": fileToken})
}

// DownloadFileHandler 下载文件
// GET /api/file/:filetoken
func DownloadFileHandler(c *gin.Context) {
	fileToken := c.Param("filetoken")
	fileReader, fileSize, err := fileSvc.DownloadFile(c.Request.Context(), fileToken)
	if err != nil {
		handleFileServiceError(c, err)
		return
	}
	defer func() { _ = fileReader.Close() }()

	c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
	c.Header("Content-Disposition", "inline; filename=\""+fileToken+"\"")
	c.DataFromReader(http.StatusOK, fileSize, "application/octet-stream", fileReader, nil)
}

func handleFileServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, fileSvc.ErrFileParamInvalid),
		errors.Is(err, fileSvc.ErrFileTokenInvalid),
		errors.Is(err, fileSvc.ErrFileHashInvalid),
		errors.Is(err, fileSvc.ErrFileSizeInvalid),
		errors.Is(err, fileSvc.ErrFileHashMismatch),
		errors.Is(err, fileSvc.ErrFileSizeMismatch),
		errors.Is(err, fileSvc.ErrFileHashSizeConflict):
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	case errors.Is(err, fileSvc.ErrFileNotFound):
		apiException.AbortWithException(c, apiException.NotFound, err)
		return
	default:
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
}
