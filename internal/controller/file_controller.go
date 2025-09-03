package controller

import (
	"net/http"

	"github.com/fikrialwan/FitByte/internal/service"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file from request: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file size (max 100KB as per your spec)
	if header.Size > 100*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 100KB limit",
		})
		return
	}

	// Validate file type (jpeg, jpg, png)
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only JPEG, JPG, and PNG are allowed",
		})
		return
	}

	// Upload to S3
	fileURL, err := c.fileService.UploadToS3(file, header.Filename, contentType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to S3: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uri": fileURL,
	})
}
