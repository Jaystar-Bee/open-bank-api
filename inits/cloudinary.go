package inits

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary
var Ctx = context.Background()

func InitCloudinary() {
	var cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}
	Cloudinary = cld
	// Upload an image to your Cloudinary product environment from a specified URL.
	//
	// Alternatively you can provide a path to a local file on your filesystem,
	// base64 encoded string, io.Reader and more.
	//
	// For additional information see:
	// https://cloudinary.com/documentation/upload_images
	//
	// Upload can be greatly customized by specifying uploader.UploadParams,
	// in this case we set the Public ID of the uploaded asset to "logo".
	// uploadResult, err := cld.Upload.Upload(
	// 	ctx,
	// 	"https://res.cloudinary.com/demo/image/upload/v1598276026/docs/models.jpg",
	// 	uploader.UploadParams{PublicID: "models",
	// 		UniqueFilename: api.Bool(false),
	// 		Overwrite:      api.Bool(true)})
	// if err != nil {
	// 	log.Fatalf("Failed to upload file, %v\n", err)
	// }
	// log.Println(uploadResult.SecureURL)

}
