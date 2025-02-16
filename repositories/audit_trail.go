package repositories

import (
	"fmt"

	database "schoolPayment/configs"
	"schoolPayment/models"
)

type AuditTrailRepository interface {
	CreateDataAuditTrail(auditTrail *models.AuditTrail) error
	InvalidateFirebaseToken(firebaseToken string, userId int) error
}

type auditTrailRepository struct{}

func NewAuditTrailRepository() AuditTrailRepository {
	return &auditTrailRepository{}
}

func (r *auditTrailRepository) CreateDataAuditTrail(auditTrail *models.AuditTrail) error {
	result := database.DB.Create(&auditTrail)
	return result.Error
}

func GetDataAuditTrailByUserId(userId int) ([]models.AuditTrail, error) {
	var auditTrails []models.AuditTrail
	// result := database.DB.Where("user_id = ?", userId).Find(auditTrails)
	query := `select * from audit_trails at`
	query += fmt.Sprintf(" where at.firebase_id != '' and at.user_id = %d and at.is_valid = true", userId)
	query += " order by at.created_at desc limit 1"
	result := database.DB.Raw(query).Scan(&auditTrails)
	if result.Error != nil {
		return []models.AuditTrail{}, result.Error
	}
	return auditTrails, nil
}

func GetAllAuditrailsNotificationMobile(schoolId int) ([]models.AuditTrail, error) {
	var auditTrails []models.AuditTrail

	query := `
        SELECT distinct at.firebase_id
        FROM audit_trails at
        Left join user_schools us on us.user_id = at.user_id
        WHERE at.firebase_id != '' AND at.platform = 'mobile' and us.school_id = ?
    `
	result := database.DB.Raw(query, schoolId).Scan(&auditTrails)
	if result.Error != nil {
		return []models.AuditTrail{}, result.Error
	}
	return auditTrails, nil
}

func (r *auditTrailRepository) InvalidateFirebaseToken(firebaseToken string, userId int) error {
	// Update the 'is_valid' column to false for matching firebaseToken and userId
	result := database.DB.Model(&models.AuditTrail{}).
		Where("firebase_id = ? AND user_id = ?", firebaseToken, userId).
		Update("is_valid", false)

	if result.Error != nil {
		return result.Error // Return error if the update failed
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no records found to invalidate for given token and userId")
	}

	return nil
}
