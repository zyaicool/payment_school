package controllers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"schoolPayment/constants"

	"github.com/gofiber/fiber/v2"
)

type AssetsController struct{}

func NewAssetsController() AssetsController {
	return AssetsController{}
}

// @Summary Get School Logo Image
// @Description Retrieve the school logo image by filename
// @Tags Assets
// @Accept json
// @Produce image/jpeg, image/png
// @Param filename path string true "Logo filename"
// @Success 200 {file} file "Image file"
// @Failure 415 {object} map[string]interface{} "Unsupported image format"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve or encode image"
// @Router /api/v1/assets/school/logo/{filename} [get]
func (a *AssetsController) GetImageLogo(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Generate the image from the filename
	img, format, err := GenerateImage("upload/school/logo/" + filename) // Adjust path as necessary
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorTextMessage + err.Error())
	}

	// Set the content type based on the format
	switch format {
	case "jpeg":
		c.Set(constants.ContentType, constants.ImageJPEG)
		err = jpeg.Encode(c.Response().BodyWriter(), img, nil)
	case "png":
		c.Set(constants.ContentType, constants.ImagePNG)
		err = png.Encode(c.Response().BodyWriter(), img)
	default:
		return c.Status(http.StatusUnsupportedMediaType).SendString(constants.UnsupportedImageFormatMessage + format)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorEncodingImageMessage + err.Error())
	}

	return nil
}

// @Summary Get Payment Method Logo Image
// @Description Retrieve the payment method logo image by filename
// @Tags Assets
// @Accept json
// @Produce image/jpeg, image/png
// @Param filename path string true "Payment method logo filename"
// @Success 200 {file} file "Image file"
// @Failure 415 {object} map[string]interface{} "Unsupported image format"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve or encode image"
// @Router /api/v1/assets/paymentMethod/logo/{filename} [get]
func (a *AssetsController) GetImagePaymentMethodLogo(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Generate the image from the filename
	img, format, err := GenerateImage("upload/paymentMethod/logo/" + filename) // Adjust path as necessary
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorTextMessage + err.Error())
	}

	// Set the content type based on the format
	switch format {
	case "jpeg":
		c.Set(constants.ContentType, constants.ImageJPEG)
		err = jpeg.Encode(c.Response().BodyWriter(), img, nil)
	case "png":
		c.Set(constants.ContentType, constants.ImagePNG)
		err = png.Encode(c.Response().BodyWriter(), img)
	default:
		return c.Status(http.StatusUnsupportedMediaType).SendString(constants.UnsupportedImageFormatMessage + format)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorEncodingImageMessage + err.Error())
	}

	return nil
}

// @Summary Get User Image
// @Description Retrieve the user image by filename
// @Tags Assets
// @Accept json
// @Produce image/jpeg, image/png
// @Param filename path string true "User image filename"
// @Success 200 {file} file "Image file"
// @Failure 415 {object} map[string]interface{} "Unsupported image format"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve or encode image"
// @Router /api/v1/assets/user/image/{filename} [get]
func (a *AssetsController) GetImageUser(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Generate the image from the filename
	img, format, err := GenerateImage("upload/user/image/" + filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorTextMessage + err.Error())
	}

	// Set the content type based on the format
	switch format {
	case "jpeg":
		c.Set(constants.ContentType, constants.ImageJPEG)
		err = jpeg.Encode(c.Response().BodyWriter(), img, nil)
	case "png":
		c.Set(constants.ContentType, constants.ImagePNG)
		err = png.Encode(c.Response().BodyWriter(), img)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorEncodingImageMessage + err.Error())
	}

	return nil
}

// @Summary Get Letterhead Image
// @Description Retrieve the school letterhead image by filename
// @Tags Assets
// @Accept json
// @Produce image/jpeg, image/png
// @Param filename path string true "Letterhead image filename"
// @Success 200 {file} file "Image file"
// @Failure 415 {object} map[string]interface{} "Unsupported image format"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve or encode image"
// @Router /api/v1/assets/school/letterhead/{filename} [get]
func (a *AssetsController) GetImageLetterhead(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Generate the image from the filename
	img, format, err := GenerateImage("upload/school/letterhead/" + filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorTextMessage + err.Error())
	}

	// Set the content type based on the format
	switch format {
	case "jpeg":
		c.Set(constants.ContentType, constants.ImageJPEG)
		err = jpeg.Encode(c.Response().BodyWriter(), img, nil)
	case "png":
		c.Set(constants.ContentType, constants.ImagePNG)
		err = png.Encode(c.Response().BodyWriter(), img)
	default:
		return c.Status(http.StatusUnsupportedMediaType).SendString(constants.UnsupportedImageFormatMessage + format)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorEncodingImageMessage + err.Error())
	}

	return nil
}

// @Summary Get Announcement Image
// @Description Retrieve the announcement image by filename
// @Tags Assets
// @Accept json
// @Produce image/jpeg, image/png
// @Param filename path string true "Announcement image filename"
// @Success 200 {file} file "Image file"
// @Failure 415 {object} map[string]interface{} "Unsupported image format"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve or encode image"
// @Router /api/v1/assets/announcement/{filename} [get]
func (a *AssetsController) GetImageAnnouncement(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Generate the image from the filename
	img, format, err := GenerateImage("upload/announcement/image/" + filename)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorTextMessage + err.Error())
	}

	// Set the content type based on the format
	switch format {
	case "jpeg":
		c.Set(constants.ContentType, constants.ImageJPEG)
		err = jpeg.Encode(c.Response().BodyWriter(), img, nil)
	case "png":
		c.Set(constants.ContentType, constants.ImagePNG)
		err = png.Encode(c.Response().BodyWriter(), img)
	default:
		return c.Status(http.StatusUnsupportedMediaType).SendString(constants.UnsupportedImageFormatMessage + format)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(constants.ErrorEncodingImageMessage + err.Error())
	}

	return nil
}

func GenerateImage(filename string) (image.Image, string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	return img, format, nil
}
