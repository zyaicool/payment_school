package controllers

import (
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

type CityController struct {
	cityService services.CityService
}

func NewCityController(cityService services.CityService) *CityController {
	return &CityController{cityService: cityService}
}

// @Summary Get Cities by Province ID
// @Description Retrieve city data filtered by Province ID
// @Tags Cities
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param provinceId query int false "Province ID"
// @Success 200 {object} map[string]interface{} "List of cities"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/city [get]
func (cityController *CityController) GetDataCities(c *fiber.Ctx) error {
	provinceID := c.QueryInt("provinceId", 0)
	districts, err := cityController.cityService.GetDataCities("data/mst_district.json", provinceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": districts,
	})
}
