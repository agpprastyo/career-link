package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// UploadAvatarFile uploads an avatar file to storage
func (r *UserRepository) UploadAvatarFile(ctx context.Context, fileName string, fileContent io.Reader, fileSize int64, contentType string) (avatarURL string, err error) {
	if fileContent == nil {
		return "", errors.New("file content cannot be nil")
	}

	// Upload the file to MinIO
	err = r.minio.UploadFile(ctx, fileName, fileContent, fileSize, contentType)
	if err != nil {
		r.log.WithError(err).WithField("fileName", fileName).Error("Failed to upload avatar file")
		return "", fmt.Errorf("failed to upload avatar file: %w", err)
	}

	// Generate URL for the uploaded file with long expiration (1 year)
	url, err := r.minio.GetFileURL(ctx, fileName, 24*7*time.Hour)
	if err != nil {
		r.log.WithError(err).WithField("fileName", fileName).Error("Failed to generate URL for avatar file")
		return "", fmt.Errorf("failed to generate URL for avatar file: %w", err)
	}

	return url, nil
}

// DeleteAvatarFile removes an avatar file from storage
func (r *UserRepository) DeleteAvatarFile(ctx context.Context, fileName string) error {
	if fileName == "" {
		return errors.New("file name cannot be empty")
	}

	// If the filename includes the "avatars/" prefix, use it as is
	// Otherwise, assume it might be a full URL and extract just the filename
	if len(fileName) > 0 && !strings.HasPrefix(fileName, "avatars/") {
		// If it's a URL, extract just the filename part
		parts := strings.Split(fileName, "/")
		if len(parts) > 0 {
			fileName = "avatars/" + parts[len(parts)-1]
		}
	}

	// Delete the file from MinIO
	err := r.minio.DeleteFile(ctx, fileName)
	if err != nil {
		r.log.WithError(err).WithField("fileName", fileName).Error("Failed to delete avatar file")
		return fmt.Errorf("failed to delete avatar file: %w", err)
	}

	r.log.WithField("fileName", fileName).Info("Successfully deleted avatar file")
	return nil
}

// Get the avatar URL for a user
func (r *UserRepository) GetAvatarURL(ctx context.Context, fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("file name cannot be empty")
	}

	// Generate URL for the uploaded file with long expiration (1 year)
	url, err := r.minio.GetFileURL(ctx, fileName, 24*7*time.Hour)
	if err != nil {
		r.log.WithError(err).WithField("fileName", fileName).Error("Failed to generate URL for avatar file")
		return "", fmt.Errorf("failed to generate URL for avatar file: %w", err)
	}

	return url, nil
}
