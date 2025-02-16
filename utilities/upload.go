package utilities

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const maxFileSize = 2 * 1024 * 1024 // 2MB
var allowedExtensions = []string{".jpg", ".jpeg", ".png"}

func UploadImage(file *multipart.FileHeader, c *fiber.Ctx, subFolder string, field string) (string, error) {
	// Set default field if empty
	if field == "" {
		field = "schoolLogo"
	}

	epoch := time.Now().Unix()

	if file.Size == 0 {
		return "", fmt.Errorf("%s cannot be empty", field)
	}

	if file.Size > maxFileSize {
		return "", fmt.Errorf("%s exceeds maximum size of 2MB", field)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	isValidExt := false
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return "", fmt.Errorf("%s must be a JPG, JPEG, or PNG file", field)
	}

	extension := filepath.Ext(file.Filename)

	newFileName := fmt.Sprintf("%d%s", epoch, extension)

	// Ensure the upload directory exists
	uploadDir := "./upload/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0o755) // Create the directory with necessary permissions
		if err != nil {
			return "", fmt.Errorf("Unable to create upload directory.")
		}
	}

	subFolderDir := SetSubFolder(subFolder)
	if _, err := os.Stat(subFolderDir); os.IsNotExist(err) {
		err = os.MkdirAll(subFolderDir, 0o755) // Create the directory with necessary permissions
		if err != nil {
			return "", fmt.Errorf("Unable to create sub folder %s", subFolder)
		}
	}

	filePath := filepath.Join(subFolderDir, newFileName)
	if err := c.SaveFile(file, filePath); err != nil {
		return "", fmt.Errorf("Unable to save image.")
	}
	cleanedPath := strings.TrimPrefix(filePath, "upload\\")
	cleanedPath = strings.TrimPrefix(cleanedPath, "upload//")
	return cleanedPath, nil
}

func SetSubFolder(subFolder string) string {
	var folderName string
	if subFolder == "school_logo" {
		folderName = "./upload/school/logo"
	}
	if subFolder == "payment_method_logo" {
		folderName = "./upload/paymentMethod/logo"
	}
	if subFolder == "school_letterhead" {
		folderName = "./upload/school/letterhead"
	}
	if subFolder == "user_image" {
		folderName = "./upload/user/image"
	}
	if subFolder == "announcement_image" {
		folderName = "./upload/announcement/image"
	}
	return folderName
}

func ConvertPath(path string) string {
	baseUrl := os.Getenv("BASE_URL")
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimPrefix(path, "upload/")
	path = baseUrl + "/api/assets/" + path
	return path
}

func ConvertPathImage(path string) *string {
	if path == "" {
		return nil
	}
	baseUrl := os.Getenv("BASE_URL")
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimPrefix(path, "upload/")
	result := baseUrl + "/api/assets/" + path
	return &result
}
