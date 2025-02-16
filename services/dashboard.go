package services

import (
	"fmt"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
	"time"
)

type DashboardService struct {
	userRepository        repositories.UserRepository
	schoolClassRepository repositories.SchoolClassRepositoryInterface
	dashboardRepository   *repositories.DashboardRepository
}

func NewDashboardService(userRepository repositories.UserRepository, schoolClassRepository repositories.SchoolClassRepositoryInterface, dashboardRepository *repositories.DashboardRepository) DashboardService {
	return DashboardService{
		userRepository:        userRepository,
		schoolClassRepository: schoolClassRepository,
		dashboardRepository:   dashboardRepository,
	}
}

func (dashboardService *DashboardService) GetDashboardFromParent(user models.User, studentID int) (*response.DashboardResponse, error) {
	var dashboardResponse response.DashboardResponse

	// Fetch students either by userID or specific studentID
	var students []models.Student
	var err error
	if studentID == 0 {
		students, err = repositories.GetStudentByUserIdDashboard(uint(user.ID))
	} else {
		var student *models.Student
		student, err = repositories.GetStudentByIDOnlyStudent(uint(studentID))
		students = append(students, *student)
	}

	if err != nil {
		return nil, err
	}

	// Process each student to populate the dashboard
	for _, student := range students {
		detailDashboardResponse := response.DetailDashboardResponse{}
		schoolClassName, schoolName := "", ""

		// Get school class and school name for the student
		if schoolClass, err := dashboardService.schoolClassRepository.GetSchoolClassByID(student.SchoolClassID); err == nil {
			schoolClassName = schoolClass.SchoolClassName
		}
		if school, err := repositories.GetSchoolByStudentId(student.ID); err == nil {
			schoolName = school.SchoolName
		}

		// Populate student response
		detailDashboardResponse.DetailDataStudent = response.StudentResponse{
			ID:          student.ID,
			Nis:         student.Nis,
			FullName:    student.FullName,
			Status:      student.Status,
			SchoolClass: schoolClassName,
			SchoolName:  schoolName,
		}

		// Fetch latest billings for the student
		if latestBillings, err := repositories.GetBillingStudentsDashboard(int(student.ID), 2); err == nil {
			for _, billing := range latestBillings {
				detailDashboardResponse.ListLatestBilling = append(detailDashboardResponse.ListLatestBilling, response.ListLatestBillingResponse{
					BillingName: billing.DetailBillingName,
					StudentName: student.FullName,
					DueDate:     billing.DueDate.Format("2006-01-02"), // Format date as string
					Amount:      billing.Amount,
					Status:      billing.PaymentStatus,
				})
			}
		} else {
			fmt.Printf("failed to fetch latest billings: %v\n", err)
		}

		// Fetch the latest donation
		if latestDonation, err := repositories.GetLatestDonation(true, student.ID, user.UserSchool.SchoolID); err == nil {
			detailDashboardResponse.DonationName = latestDonation.BillingName
		} else {
			fmt.Printf("failed to fetch latest donation: %v\n", err)
		}

		// Add to the response
		dashboardResponse.DataDashboard = append(dashboardResponse.DataDashboard, detailDashboardResponse)
	}

	return &dashboardResponse, nil
}

func (dashboardService *DashboardService) GetDashboardFromAdmin(userID int) (*response.DashboardAdminResponse, error) {
	var dashboardResponse response.DashboardAdminResponse

	// Get user's school ID
	user, err := dashboardService.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// Get active students count
	activeStudents, err := dashboardService.dashboardRepository.GetActiveStudentsCount(user.UserSchool.SchoolID)
	if err != nil {
		return nil, err
	}
	dashboardResponse.ActiveStudent = activeStudents

	// Get total transactions for last month
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := now
	totalTransactions, err := dashboardService.dashboardRepository.GetTotalTransactionsLastMonth(user.UserSchool.SchoolID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	dashboardResponse.TotalTransaction = totalTransactions

	// Get total billings for last month
	totalBillings, err := dashboardService.dashboardRepository.GetTotalBillingsLastMonth(user.UserSchool.SchoolID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	dashboardResponse.TotalBilling = totalBillings

	return &dashboardResponse, nil
}
