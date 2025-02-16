package repositories

import (
	database "schoolPayment/configs"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type BankAccountRepository interface {
	CreateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error)
	CheckAccountNumberExists(accountNumber string) (bool, error)
	GetBankAccountByID(id uint) (*models.BankAccount, error)
	UpdateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error)
	CheckAccountNumberExistsExcept(accountNumber string, bankID uint) (bool, error)
	GetAllBankAccounts(page, limit int, search string, sortBy string, sortOrder string, user models.User) ([]models.BankAccount, int, error)
	GetUserById(ids []int, users *[]models.User) error
	DeleteBankAccount(id uint) error
	FindUsingBankAccountOnBilling(bankAccountId int) bool
}

type bankAccountRepository struct{
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepository{db: db}
}

func (r *bankAccountRepository) CreateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	result := r.db.Create(&bankAccount)
	return bankAccount, result.Error
}

// CheckAccountNumberExists checks if an account number already exists
func (r *bankAccountRepository) CheckAccountNumberExists(accountNumber string) (bool, error) {
	var count int64
	err := database.DB.Model(&models.BankAccount{}).Where("account_number = ?", accountNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckAccountNumberExistsExcept checks if an account number exists for another bank account
func (r *bankAccountRepository) CheckAccountNumberExistsExcept(accountNumber string, bankID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.BankAccount{}).
		Where("account_number = ? AND id != ?", accountNumber, bankID).
		Count(&count).Error
	return count > 0, err
}

// GetBankAccountByID retrieves a bank account by its ID
func (r *bankAccountRepository) GetBankAccountByID(id uint) (*models.BankAccount, error) {
	var bankAccount models.BankAccount
	err := database.DB.Preload("School").First(&bankAccount, id).Error
	if err != nil {
		return nil, err
	}
	return &bankAccount, nil
}

// UpdateBankAccount updates an existing bank account in the database
func (r *bankAccountRepository) UpdateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	result := r.db.Save(&bankAccount)
	return bankAccount, result.Error
}

func (r *bankAccountRepository) GetAllBankAccounts(page, limit int, search string, sortBy string, sortOrder string, user models.User) ([]models.BankAccount, int, error) {
	offset := (page - 1) * limit

	// Prepare the query for bank accounts
	query := database.DB.Model(&models.BankAccount{}).Preload("School")

	if search != "" {
		query = query.Where("bank_accounts.account_name ILIKE ? OR bank_accounts.bank_name ILIKE ? OR bank_accounts.account_number ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if user.UserSchool != nil {
		query = query.Joins("JOIN schools ON schools.id = bank_accounts.school_id ").
			Where("schools.id = ?", user.UserSchool.School.ID)
	}

	// Fetch bank accounts
	var bankAccounts []models.BankAccount
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("CASE WHEN bank_accounts.updated_at IS NOT NULL THEN 0 ELSE 1 END, bank_accounts.updated_at DESC, bank_accounts.created_at DESC")
	}

	if err := query.Find(&bankAccounts).Error; err != nil {
		return nil, 0, err
	}

	// Get total count after applying the same conditions
	var totalData int64
	if user.UserSchool != nil {
		if err := database.DB.Model(&models.BankAccount{}).Joins("JOIN schools ON schools.id = bank_accounts.school_id ").
			Where("schools.id = ?", user.UserSchool.School.ID).
			Where("account_name ILIKE ? OR bank_name ILIKE ? OR account_number ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%").
			Count(&totalData).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := database.DB.Model(&models.BankAccount{}).
			Where("account_name ILIKE ? OR bank_name ILIKE ? OR account_number ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%").
			Count(&totalData).Error; err != nil {
			return nil, 0, err
		}
	}

	return bankAccounts, int(totalData), nil
}

func (r *bankAccountRepository) GetUserById(ids []int, users *[]models.User) error {
	if len(ids) == 0 {
		return nil // No IDs to fetch
	}
	return database.DB.Where("id IN ?", ids).Find(users).Error
}

// DeleteBankAccount deletes a bank account by ID
func (r *bankAccountRepository) DeleteBankAccount(id uint) error {
	if err := r.db.Delete(&models.BankAccount{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *bankAccountRepository) FindUsingBankAccountOnBilling(bankAccountId int) bool {
	var billing []models.Billing
	_ = database.DB.Where("bank_account_id = ?", bankAccountId).Find(&billing)

	if len(billing) > 0 {
		return false
	}

	return true
}
