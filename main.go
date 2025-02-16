package main

import (
	"log"
	"os"

	"schoolPayment/configs"
	config "schoolPayment/configs"
	"schoolPayment/controllers"
	"schoolPayment/repositories"
	routes "schoolPayment/routes"
	"schoolPayment/services"

	_ "schoolPayment/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
)

type SwaggerDoc struct {
	Servers []map[string]string `json:"servers"`
}

func main() {
	// Load environment variables
	config.LoadEnvVariables()

	// Connect to the database
	config.ConnectToDatabase()

	// Run Liquibase migration
	if err := config.RunLiquibaseUpdate(); err != nil {
		log.Fatalf("Liquibase migration failed: %v", err)
	}

	// Setup Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Specify your frontend local development URL
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: true,
	}))

	// Initialize Repositories
	userRepository := repositories.NewUserRepository()
	roleRepository := repositories.NewRoleRepository()
	studentRepository := repositories.NewStudentRepository(config.DB)
	studentParentRepository := repositories.NewStudentParentRepository()
	transactionRepository := repositories.NewTransactionRepository()
	schoolYearRepository := repositories.NewSchoolYearRepository(configs.DB)
	schoolGradeRepository := repositories.NewSchoolGradeRepository(config.DB)
	schoolRepository := repositories.NewSchoolRepository(configs.DB)
	schoolClassRepository := repositories.NewSchoolClassRepository()
	billingHistoryRepository := repositories.NewBillingHistoryRepository()
	billingRepository := repositories.NewBillingRepository(config.DB)
	billingStudentRepository := repositories.NewBillingStudentRepository(config.DB)
	bankAccountRepository := repositories.NewBankAccountRepository(config.DB)
	schoolMajorRepository := repositories.NewSchoolMajorRepository(config.DB)
	paymentMethodRepository := repositories.NewPaymentMethodRepository(configs.DB)
	auditTrailRepository := repositories.NewAuditTrailRepository()
	scheduleRepository := repositories.NewScheduleRepository()
	paymentReportRepository := repositories.NewPaymentReportRepository()
	billingReportRepository := repositories.NewBillingReportRepository()
	dashboardRepository := repositories.NewDashboardRepository(configs.DB)
	announcementRepository := repositories.NewAnnouncementRepository(configs.DB)
	invoiceFormatRepository := repositories.NewInvoiceFormatRepository(configs.DB)

	// Initialize Services
	userService := services.NewUserService(userRepository, roleRepository, schoolRepository)
	roleService := services.NewRoleService(roleRepository)
	studentParentService := services.NewStudentParentService(studentParentRepository, userRepository)
	billingHistoryService := services.NewBillingHistoryService(billingHistoryRepository, userRepository, schoolRepository, paymentMethodRepository)
	transactionService := services.NewTransactionService(transactionRepository, userRepository, schoolRepository, billingHistoryService, billingStudentRepository, billingRepository, studentRepository)
	schoolYearService := services.NewSchoolYearService(schoolYearRepository, userRepository)
	schoolGradeService := services.NewSchoolGradeService(schoolGradeRepository)
	schoolService := services.NewSchoolService(schoolRepository, userRepository)
	billingService := services.NewBillingService(billingRepository, userRepository, schoolClassRepository, schoolYearRepository, schoolGradeRepository, studentRepository)
	billingStudentService := services.NewBillingStudentService(billingStudentRepository, userRepository, schoolYearRepository, schoolClassRepository, billingRepository, schoolGradeRepository, studentRepository)
	bankAccountService := services.NewBankAccountService(bankAccountRepository, userRepository)
	loginService := services.NewLoginService(userRepository, auditTrailRepository)
	dashboardService := services.NewDashboardService(userRepository, schoolClassRepository, dashboardRepository)
	schoolMajorService := services.NewSchoolMajorService(schoolMajorRepository, userRepository)
	paymentMethodService := services.NewPaymentMethodService(paymentMethodRepository)
	scheduleService := services.NewScheduleService(scheduleRepository, userRepository, billingRepository)
	paymentReportService := services.NewPaymentReportService(paymentReportRepository, userRepository)
	billingReportService := services.NewBillingReportService(billingReportRepository, userRepository)
	announcementService := services.NewAnnouncementService(announcementRepository, userRepository)
	invoiceFormatService := services.NewInvoiceFormatService(invoiceFormatRepository)

	// Initialize Controllers
	userController := controllers.NewUserController(userService)
	roleController := controllers.NewRoleController(roleService)
	studentParentController := controllers.NewStudentParentController(studentParentService)
	transactionController := controllers.NewTransactionController(transactionService)
	schoolYearController := controllers.NewSchoolYearController(schoolYearService)
	schoolGradeController := controllers.NewSchoolGradeController(schoolGradeService)
	schoolController := controllers.NewSchoolController(schoolService)
	bilingHistoryController := controllers.NewBillingHistoryController(billingHistoryService)
	billingController := controllers.NewBillingController(billingService)
	billingStudentController := controllers.NewBillingStudentController(billingStudentService)
	bankAccountController := controllers.NewBankAccountController(bankAccountService)
	loginController := controllers.NewLoginController(loginService)
	dashController := controllers.NewDashboardController(dashboardService)
	schoolMajorController := controllers.NewSchoolMajorController(schoolMajorService)
	assetsController := controllers.NewAssetsController()
	paymentMethodController := controllers.NewPaymentMethodController(paymentMethodService)
	scheduleController := controllers.NewScheduleController(scheduleService)
	paymentReportController := controllers.NewPaymentReportController(paymentReportService)
	billingReportController := controllers.NewBillingReportController(billingReportService)
	announcementController := controllers.NewAnnouncementController(announcementService)
	invoiceFormatController := controllers.NewInvoiceFormatController(invoiceFormatService)

	// Setup routes
	api := app.Group("/v1")
	routes.SetupRoleRoutes(api, roleController)
	routes.SetupUserRoutes(api, userController)
	routes.SetupParentRoutes(api, studentParentController)
	routes.SetupStudentRoutes(api)
	routes.SetupBillingRoutes(api, billingController)
	routes.SetupSchoolYearRoutes(api, schoolYearController)
	routes.SetupSchoolRoutes(api, schoolController)
	routes.SetupSchoolGradeRoutes(api, schoolGradeController)
	routes.SetupSchoolClassRoutes(api)
	routes.SetupBillingStudentRoutes(api, billingStudentController)
	routes.SetupTransactionRoutes(api, transactionController)
	routes.SetupBillingHistoryRoutes(api, bilingHistoryController)
	routes.SetupDashboardRoutes(api, dashController)
	routes.SetupBankAccountRoutes(api, &bankAccountController)
	routes.SetupPrefixClassRoutes(api)
	routes.SetupSchoolMajorRoutes(api, schoolMajorController.(*controllers.SchoolMajorController))
	routes.SetupPaymentMethodRoutes(api, &paymentMethodController)
	routes.SetupScheduleRoutes(api, scheduleController)
	routes.SetupLoginRoutes(api, loginController)
	routes.SetupAssetsRoutes(app, &assetsController)
	routes.SetupPaymentReportRoutes(api, paymentReportController)
	routes.SetupBillingReportRoutes(api, billingReportController)
	routes.SetupAnnouncementRoutes(api, announcementController)
	routes.SetupInvoiceFormatRoutes(api, invoiceFormatController)
	routes.SetupRoutes(api)

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port)) // Ensure single colon here
}
