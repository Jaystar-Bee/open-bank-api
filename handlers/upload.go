package handlers

import (
	"fmt"
	"net/http"

	"github.com/Jaystar-Bee/open-bank-api/inits"
	"github.com/Jaystar-Bee/open-bank-api/models"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// UploadFile godoc
//	@Summary		Upload File
//	@Description	Upload File
//	@Tags			Upload
//	@Accept			multipart/form-data
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"Upload File"
//	@Produce		json
//	@Success		200	{object}	models.HTTP_FILE_RESPONSE
//	@Failure		400	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/upload [post]
func UploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
		})
		return
	}
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
		})
		return
	}
	userId := context.GetInt64("user")

	user, err := models.GetUserByID(userId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to get user",
		})
		return
	}

	res, err := inits.Cloudinary.Upload.Upload(inits.Ctx, file, uploader.UploadParams{
		Overwrite:      api.Bool(true),
		UseFilename:    api.Bool(true),
		UniqueFilename: api.Bool(false),
		Folder:         fmt.Sprintf("OpenBank/%s%s-%d", user.FirstName, user.LastName, user.ID),
		AssetFolder:    fmt.Sprint(user.ID),
		Tags:           []string{"profile"},
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
