package upload

import (
	"context"
	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"time"
)

func ImageUploadHelper(input interface{}) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cloudName := config.InitConfig().CCName
	apiKey := config.InitConfig().CCAPIKey
	apiSecret := config.InitConfig().CCAPISecret
	uploadFolder := config.InitConfig().CCFolder

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		Folder: uploadFolder,
	}

	uploadResult, err := cld.Upload.Upload(ctx, input, uploadParams)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
