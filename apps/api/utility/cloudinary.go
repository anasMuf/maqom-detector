package utility

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(fileBytes []byte, filename string) (string, error) {
	cldURL := os.Getenv("CLOUDINARY_URL")
	if cldURL == "" {
		return "", nil // Or handle error
	}

	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	
	// Use bytes.NewReader for better SDK compatibility
	fileReader := bytes.NewReader(fileBytes)

	uploadResult, err := cld.Upload.Upload(ctx, fileReader, uploader.UploadParams{
		PublicID:     filename,
		ResourceType: "video", // Audio is treated as video in Cloudinary
		Folder:       "maqom-detector/audio",
	})
	
	if err != nil {
		fmt.Printf("Cloudinary Error: %v\n", err)
		return "", err
	}

	fmt.Printf("Cloudinary Success: %s\n", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}
