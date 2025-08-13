package cloudinary

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Client is a struct for talking to the Cloudinary API.
type Client struct {
	cld *cloudinary.Cloudinary
}

// NewClient creates a new Cloudinary client.
// The cloudName, apiKey, and apiSecret come from your Cloudinary dashboard.
func NewClient(cloudName, apiKey, apiSecret string) (*Client, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary client: %w", err)
	}

	return &Client{cld: cld}, nil
}

// UploadImage takes image data as bytes and uploads it to Cloudinary.
func (c *Client) UploadImage(imageData []byte) (string, error) {
	// The context object can be used to set deadlines or cancel requests.
	// For this simple case, context.Background() is sufficient.
	ctx := context.Background()
	reader := bytes.NewReader(imageData)
	// The SDK's upload function takes the image data (as a reader) and parameters.
	// We use an empty uploader.UploadParams{} for a basic upload.
	uploadResult, err := c.cld.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to cloudinary: %w", err)
	}

	// The result contains a SecureURL, which is the HTTPS link to the image.
	return uploadResult.SecureURL, nil
}
