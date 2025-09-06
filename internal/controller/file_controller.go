package controller

import (
	"net/http"
	"strings"

	"github.com/fikrialwan/FitByte/internal/service"
	"github.com/fikrialwan/FitByte/pkg/handler"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) FileController {
	return FileController{
		fileService: fileService,
	}
}

// UploadFile godoc
// @Summary Upload file to S3
// @Description Upload an image file (JPEG, JPG, PNG) to S3 storage with max size of 100KB
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload (max 100KB, JPEG/JPG/PNG only)"
// @Success 200 {object} map[string]string "Returns S3 file URL"
// @Failure 400 {object} map[string]string "Bad request - invalid file or size"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /file [post]
func (c FileController) UploadFile(ctx *gin.Context) {
	// Get file from form data
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		handler.ResponseError(ctx, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// Validate file size (max 100KB as per your spec)
	if header.Size > 100*1024 {
		handler.ResponseError(ctx, http.StatusBadRequest, "File size exceeds 100KB limit")
		return
	}

	// Validate file type by extension (more reliable than content-type for multipart uploads)
	contentType := header.Header.Get("Content-Type")
	filename := header.Filename

	// Check file extension
	validExtensions := []string{".jpg", ".jpeg", ".png"}
	validExtension := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			validExtension = true
			break
		}
	}

	if !validExtension {
		handler.ResponseError(ctx, http.StatusBadRequest, "Invalid file type. Only JPEG, JPG, and PNG are allowed")
		return
	}

	// Upload to S3
	fileURL, err := c.fileService.UploadToS3(file, header.Filename, contentType)
	if err != nil {
		handler.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	handler.ResponseSuccess(ctx, http.StatusOK, gin.H{
		"uri": fileURL,
	})
}
