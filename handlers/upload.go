package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jaystar-Bee/open-bank-api/inits"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// UploadFile godoc
// @Summary Upload File
// @Description Upload File
// @Tags Upload
// @Accept multipart/form-data
// @Param file formData file true "Upload File"
// @Param id formData int true "Upload File"
// @Produce json
// @Success 200 {object} models.HTTP_FILE_RESPONSE
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /upload [post]
func UploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
		})
		return
	}
	id, err := strconv.Atoi(context.PostForm("id"))

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
		})
		return
	}
	res, err := inits.Cloudinary.Upload.Upload(inits.Ctx, file, uploader.UploadParams{
		Overwrite:      api.Bool(true),
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(false),
		Folder:         fmt.Sprint(id),
	})

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message":    "Unable to upload file",
			"dev_reason": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"data":    res.SecureURL,
	})

}
