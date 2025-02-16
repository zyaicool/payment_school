package repositories

import (
	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	GetAllPaymentMethod(search string) ([]models.PaymentMethod, error)
	GetPaymentMethodByID(id int) (*models.PaymentMethod, error)
	CreatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error)
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (paymentMethodRepository *paymentMethodRepository) GetAllPaymentMethod(search string) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod

	// Build the query
	query := paymentMethodRepository.db.Model(&models.PaymentMethod{})

	// Apply search filter if the search parameter is not empty
	if search != "" {
		query = query.Where("bank_name LIKE ?", "%"+search+"%")
	}

	// Execute the query
	result := query.Find(&paymentMethods)
	if result.Error != nil {
		return nil, result.Error
	}

	return paymentMethods, nil
}

func (paymentMethodRepository *paymentMethodRepository) GetPaymentMethodByID(id int) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	result := paymentMethodRepository.db.First(&paymentMethod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentMethod, nil
}

func (paymentMethodRepository *paymentMethodRepository) CreatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	result := paymentMethodRepository.db.Create(&paymentMethod)
	return paymentMethod, result.Error
}

func (paymentMethodRepository *paymentMethodRepository) UpdatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	result := paymentMethodRepository.db.Save(&paymentMethod)
	if result.Error != nil {
		return nil, result.Error
	}
	return paymentMethod, nil
}

func GetPaymentMethodByID(id int) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	result := database.DB.First(&paymentMethod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &paymentMethod, nil
}

func GetPaymentMethodForPaymentType() ([]response.PaymentType, error) {
	var paymentTypes []response.PaymentType

	query := `select mpm.id, concat(mpm.payment_method, mpm.bank_code) as code, concat(mpm.payment_method, ' - ', mpm.bank_name) as name
	from master_payment_method mpm`
	result := database.DB.Raw(query).Scan(&paymentTypes)
	if result.Error != nil {
		return []response.PaymentType{}, result.Error
	}
	return paymentTypes, nil
}
