package services

import (
	"errors"
	"fmt"
	"math"

	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
	"schoolPayment/utilities"
)
type BankAccountServiceInterface interface {
	CreateBankAccount(bankRequest *request.BankAccountCreateRequest, userID uint) (*models.BankAccount, error)
	UpdateBankAccount(bankID uint, bankRequest *request.BankAccountCreateRequest, userID int) (*models.BankAccount, error)
	GetAllBankAccounts(page int, limit int, search string, sortBy string, sortOrder string, user models.User) (response.BankAccountListResponse, error)
	GetBankAccountDetails(id uint) (*response.BankAccountData, error)
	DeleteBankAccount(id uint) error
}

type BankAccountService struct {
	bankAccountRepository repositories.BankAccountRepository
	userRepository        repositories.UserRepository
}

func NewBankAccountService(bankAccountRepository repositories.BankAccountRepository, userRepository repositories.UserRepository) BankAccountServiceInterface {
	return &BankAccountService{bankAccountRepository: bankAccountRepository, userRepository: userRepository}
}

func (bankAccountService *BankAccountService) CreateBankAccount(bankRequest *request.BankAccountCreateRequest, userID uint) (*models.BankAccount, error) {
	// Check if the account number already exists
	exists, err := bankAccountService.bankAccountRepository.CheckAccountNumberExists(bankRequest.AccountNumber)
	if err != nil {
		return nil, err // Return the error if there's an issue checking
	}
	if exists {
		return nil, fmt.Errorf("account number already exists")
	}

	bankAccount := &models.BankAccount{
		SchoolID:      bankRequest.SchoolID,
		BankName:      bankRequest.BankName,
		AccountName:   bankRequest.AccountName,
		AccountNumber: bankRequest.AccountNumber,
		AccountOwner:  bankRequest.AccountOwner,
		Master: models.Master{
			CreatedBy: int(userID),
		},
	}

	return bankAccountService.bankAccountRepository.CreateBankAccount(bankAccount)
}

func (bankAccountService *BankAccountService) UpdateBankAccount(bankID uint, bankRequest *request.BankAccountCreateRequest, userID int) (*models.BankAccount, error) {
	// Check if the account number already exists (excluding the current bankID)
	exists, err := bankAccountService.bankAccountRepository.CheckAccountNumberExistsExcept(bankRequest.AccountNumber, bankID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("account number already exists")
	}

	// Retrieve the existing bank account record
	bankAccount, err := bankAccountService.bankAccountRepository.GetBankAccountByID(bankID)
	if err != nil {
		return nil, fmt.Errorf("bank account not found")
	}

	// Update fields only if values are provided
	if bankRequest.SchoolID != 0 {
		bankAccount.SchoolID = bankRequest.SchoolID
	}
	if bankRequest.BankName != "" {
		bankAccount.BankName = bankRequest.BankName
	}
	if bankRequest.AccountName != "" {
		bankAccount.AccountName = bankRequest.AccountName
	}
	if bankRequest.AccountNumber != "" {
		bankAccount.AccountNumber = bankRequest.AccountNumber
	}
	if bankRequest.AccountOwner != "" {
		bankAccount.AccountOwner = bankRequest.AccountOwner
	}
	bankAccount.UpdatedBy = userID

	// Save updated bank account
	return bankAccountService.bankAccountRepository.UpdateBankAccount(bankAccount)
}

func (bankAccountService *BankAccountService) GetAllBankAccounts(page int, limit int, search string, sortBy string, sortOrder string, user models.User) (response.BankAccountListResponse, error) {
	var mapBankAccount response.BankAccountListResponse
	mapBankAccount.Limit = limit
	mapBankAccount.Page = page
	mapBankAccount.Data = []response.BankAccountData{}

	if sortBy != "" {
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	// Get the list of bank accounts and total data count
	listBankAccounts, totalData, err := bankAccountService.bankAccountRepository.GetAllBankAccounts(page, limit, search, sortBy, sortOrder, user)
	if err != nil {
		return mapBankAccount, err // Return empty response in case of an error
	}

	// Set total data and calculate total pages
	mapBankAccount.TotalData = totalData
	if limit > 0 {
		mapBankAccount.TotalPage = (totalData + limit - 1) / limit // This calculates total pages
	} else {
		mapBankAccount.TotalPage = 1 // This calculates total pages
	}

	// Calculate total pages
	if limit > 0 {
		mapBankAccount.TotalPage = int(math.Ceil(float64(totalData) / float64(limit))) // Total pages based on data and limit
	} else {
		mapBankAccount.TotalPage = 1 // If limit is 0, set to 1 to avoid division by zero
	}

	// Create a map to hold user names for fast lookup
	userIDs := make([]int, 0)
	userMap := make(map[int]string) // Maps user ID to username

	for _, account := range listBankAccounts {
		userIDs = append(userIDs, account.CreatedBy)
	}

	// Fetch user details for the user IDs
	var users []models.User
	if err := bankAccountService.bankAccountRepository.GetUserById(userIDs, &users); err == nil {
		// Map the user data to the userMap for quick lookup
		for _, user := range users {
			userMap[int(user.ID)] = user.Username
		}
	}

	// Map your BankAccount slice to BankAccountData slice
	for _, account := range listBankAccounts {
		createdByName := userMap[account.CreatedBy] // Lookup the createdBy username
		bankName, _ := GetBankName(account.BankName)
		placeHolder := fmt.Sprintf("%s - %s - %s", bankName, account.AccountName, account.AccountNumber)

		statusIsDelete := bankAccountService.bankAccountRepository.FindUsingBankAccountOnBilling(int(account.ID))
		mapBankAccount.Data = append(mapBankAccount.Data, response.BankAccountData{
			ID:            int(account.ID),
			BankName:      bankName,
			AccountNumber: account.AccountNumber,
			AccountName:   account.AccountName,
			AccountOwner:  account.AccountOwner,
			CreatedBy:     createdByName,
			CreatedAt:     account.CreatedAt,
			UpdatedBy:     account.UpdatedBy,
			UpdatedAt:     account.UpdatedAt,
			School: response.SchoolData{
				ID:   int(account.School.ID),
				Name: account.School.SchoolName,
				// Map other fields if necessary
			},
			PlaceHolder: placeHolder,
			IsDelete:    statusIsDelete, // nanti data ini ada pengecekannya
		})
	}

	return mapBankAccount, nil
}

func (bankAccountService *BankAccountService) GetBankAccountDetails(id uint) (*response.BankAccountData, error) {
	// Fetch the bank account by ID
	bankAccount, err := bankAccountService.bankAccountRepository.GetBankAccountByID(id)
	if err != nil {
		return nil, err
	}

	// Use GetUserById with a slice containing just the CreatedBy ID
	var users []models.User
	err = bankAccountService.bankAccountRepository.GetUserById([]int{int(bankAccount.CreatedBy)}, &users)
	if err != nil {
		return nil, err
	}

	// Check if user was found; if not, return an error
	var createByUser models.User
	if len(users) > 0 {
		createByUser = users[0]
	} else {
		createByUser = models.User{Username: ""}
	}

	bankName, _ := GetBankName(bankAccount.BankName)
	placeHolder := fmt.Sprintf("%s-%s-%s", bankName, bankAccount.AccountName, bankAccount.AccountNumber)

	// Map to response.BankAccountData
	bankAccountData := &response.BankAccountData{
		ID:            int(bankAccount.ID),
		BankName:      bankName,
		AccountName:   bankAccount.AccountName,
		AccountNumber: bankAccount.AccountNumber,
		AccountOwner:  bankAccount.AccountOwner,
		CreatedBy:     createByUser.Username,
		School: response.SchoolData{
			ID:   int(bankAccount.School.ID),
			Name: bankAccount.School.SchoolName,
			// Map other fields if necessary
		},
		PlaceHolder: placeHolder,
	}

	return bankAccountData, nil
}

// DeleteBankAccount deletes a bank account by ID
func (bankAccountService *BankAccountService) DeleteBankAccount(id uint) error {
	bankAccount, err := bankAccountService.bankAccountRepository.GetBankAccountByID(id)
	if err != nil {
		return err
	}

	if bankAccount == nil {
		return errors.New("bank account not found")
	}

	return bankAccountService.bankAccountRepository.DeleteBankAccount(id) // Implement this method in your repository
}

func GetBankName(bankCode string) (string, error) {
	banks, err := models.LoadBanks("data/banks.json") // Path to your JSON file
	if err != nil {
		return bankCode, nil
	}

	if len(bankCode) > 3 {
		return bankCode, nil
	}

	bankName := bankCode
	for _, bank := range banks {
		if bank.Code == bankCode {
			bankName = bank.Alias
			break
		}
	}

	return bankName, nil
}
