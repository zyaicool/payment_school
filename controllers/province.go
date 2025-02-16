package controllers

import (
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

type ProvinceController struct {
	provinceService services.ProvinceService
}

func NewProvinceController(provinceService services.ProvinceService) *ProvinceController {
	return &ProvinceController{provinceService: provinceService}
}

// @Summary Get Provinces Data
// @Description Retrieve a list of provinces from the provided data file
// @Tags Province
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of provinces"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve provinces data"
// @Router /api/v1/province [get]
func (provinceController *ProvinceController) GetDataProvinces(c *fiber.Ctx) error {
	regions, err := provinceController.provinceService.GetDataProvinces("data/mst_province.json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": regions,
	})
}
