package routes

import (
	controllers "schoolPayment/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAssetsRoutes(api fiber.Router, assetsController *controllers.AssetsController) {
	apiAssets := api.Group("/assets")

	apiAssetsUser := apiAssets.Group("/user")
	apiAssetsUser.Get("/image/:filename", assetsController.GetImageUser)

	apiAssetsSchool := apiAssets.Group("/school")
	apiAssetsSchool.Get("/logo/:filename", assetsController.GetImageLogo)
	apiAssetsSchool.Get("/letterhead/:filename", assetsController.GetImageLetterhead)

	apiAssetsPaymentMethod := apiAssets.Group("/paymentMethod")
	apiAssetsPaymentMethod.Get("/logo/:filename", assetsController.GetImagePaymentMethodLogo)

	apiAssetsAnnouncement := apiAssets.Group("/announcement")
	apiAssetsAnnouncement.Get("/image/:filename", assetsController.GetImageAnnouncement)
}
