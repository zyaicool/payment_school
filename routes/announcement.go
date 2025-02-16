package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupAnnouncementRoutes(api fiber.Router, AnnouncementController *controllers.AnnouncementController) {
	apiAnnouncement := api.Group("/announcement")
	apiAnnouncement.Get("/type", utilities.JWTProtected, AnnouncementController.GetAnnouncementTypes)
	apiAnnouncement.Post("/create", utilities.JWTProtected, AnnouncementController.CreateAnnouncement)
	apiAnnouncement.Put("/update/:id", utilities.JWTProtected, AnnouncementController.UpdateAnnouncement)
	apiAnnouncement.Delete("/delete/:id", utilities.JWTProtected, AnnouncementController.DeleteAnnouncement)
	apiAnnouncement.Get("/getList", utilities.JWTProtected, AnnouncementController.GetAnnouncementList)
	apiAnnouncement.Get("/detail/:id", utilities.JWTProtected, AnnouncementController.GetAnnouncementDetail)
}
