package repositories

import (
	"fmt"
	"schoolPayment/constants"
	"schoolPayment/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{
		db: db,
	}
}

func (r *DashboardRepository) GetActiveStudentsCount(schoolID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Student{}).
		Joins(constants.JoinUserStudentsToStudentsDashboard).
		Joins(constants.JoinUserSchoolsToUserStudents).
		Where("user_schools.school_id = ? AND students.status = ?", schoolID, "aktif").
		Count(&count).Error
	return int(count), err
}

func (r *DashboardRepository) GetTotalTransactionsLastMonth(schoolID uint, startDate, endDate time.Time) (int, error) {
	var count int64
	fmt.Println(startDate)
	fmt.Println(endDate)
	err := r.db.Model(&models.TransactionBilling{}).
		Joins("JOIN students ON students.id = transaction_billings.student_id").
		Joins(constants.JoinUserStudentsToStudentsDashboard).
		Joins(constants.JoinUserSchoolsToUserStudents).
		Where("user_schools.school_id = ? AND (transaction_billings.created_at BETWEEN ? AND ? OR transaction_billings.updated_at BETWEEN ? AND ?)",
			schoolID, startDate, endDate, startDate, endDate).
		Where("transaction_billings.deleted_at IS NULL").
		Count(&count).Error
	return int(count), err
}

func (r *DashboardRepository) GetTotalBillingsLastMonth(schoolID uint, startDate, endDate time.Time) (int, error) {
	var count int64
	err := r.db.Model(&models.TransactionBilling{}).
		Joins("JOIN students ON students.id = transaction_billings.student_id").
		Joins(constants.JoinUserStudentsToStudentsDashboard).
		Joins(constants.JoinUserSchoolsToUserStudents).
		Where("transaction_billings.transaction_status = ? AND user_schools.school_id = ? AND transaction_billings.created_at BETWEEN ? AND ?",
			"PS02", schoolID, startDate, endDate).
		Count(&count).Error
	return int(count), err
}
