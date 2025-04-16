package minio

import (
	"context"
	"fmt"
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/pkg/logger"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the MinIO client with additional functionality
type Client struct {
	minioClient *minio.Client
	bucketName  string
	location    string
	log         *logger.Logger
}

// NewClient creates a new MinIO client
func NewClient(cfg *config.AppConfig, log *logger.Logger) (*Client, error) {
	// Create MinIO client
	client, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create instance of our client wrapper
	minioClient := &Client{
		minioClient: client,
		bucketName:  cfg.Minio.BucketName,
		location:    cfg.Minio.Location,
		log:         log,
	}

	// Ensure bucket exists
	err = minioClient.EnsureBucketExists(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return minioClient, nil
}

// EnsureBucketExists checks if the configured bucket exists and creates it if not
func (c *Client) EnsureBucketExists(ctx context.Context) error {
	exists, err := c.minioClient.BucketExists(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = c.minioClient.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{
			Region: c.location,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		c.log.Infof("Bucket %s created successfully", c.bucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (c *Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.minioClient.PutObject(ctx, c.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// GetFileURL gets a resigned URL for an object with specified expiration
func (c *Client) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Get resigned URL for object
	url, err := c.minioClient.PresignedGetObject(ctx, c.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// DownloadFile downloads a file from MinIO
func (c *Client) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := c.minioClient.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return obj, nil
}

// DeleteFile deletes a file from MinIO
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := c.minioClient.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// ListFiles lists all files in a directory (prefix)
func (c *Client) ListFiles(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo
	objectCh := c.minioClient.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// Close closes the MinIO client
func (c *Client) Close() error {
	c.log.Info("Closing MinIO client")
	return nil
}
