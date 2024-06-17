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
	cld.Config.URL.Secure = false
	Cloudinary = cld

}
