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
func NewClient(cloudName, apiKey, apiSecret string) (*Client, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary client: %w", err)
	}

	return &Client{cld: cld}, nil
}

// UploadImage takes image data as bytes and uploads it to Cloudinary.
func (c *Client) UploadImage(imageData []byte) (string, error) {

	ctx := context.Background()
	reader := bytes.NewReader(imageData)

	uploadResult, err := c.cld.Upload.Upload(ctx, reader, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}
