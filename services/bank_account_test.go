package services

import (
	"fmt"
	"testing"

	"schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"

	"github.com/stretchr/testify/assert"
)

type MockBankAccountRepository struct {
	exists          bool
	errExists       error
	bankAccount     *models.BankAccount
	errCreate       error
	bankAccountByID *models.BankAccount
	errGetByID      error
	errUpdate       error
	bankAccountList []models.BankAccount
	page            int
	errGetAll       error
	totalCount      int
	errGetUser      error
	errorDelete     error
	isAccountExist  bool
}

func (m *MockBankAccountRepository) CheckAccountNumberExists(accountNumber string) (bool, error) {
	return m.exists, m.errExists
}

func (m *MockBankAccountRepository) CheckAccountNumberExistsExcept(accountNumber string, id uint) (bool, error) {
	return m.exists, m.errExists
}

func (m *MockBankAccountRepository) CreateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	if m.errCreate != nil {
		return nil, m.errCreate
	}
	return bankAccount, nil
}

func (m *MockBankAccountRepository) GetBankAccountByID(id uint) (*models.BankAccount, error) {
	return m.bankAccountByID, m.errGetByID
}

func (m *MockBankAccountRepository) UpdateBankAccount(bankAccount *models.BankAccount) (*models.BankAccount, error) {
	if m.errUpdate != nil {
		return nil, m.errUpdate
	}
	return bankAccount, nil
}

// GetAllBankAccounts simulates retrieving all bank accounts with pagination and search.
func (m *MockBankAccountRepository) GetAllBankAccounts(page, limit int, search, sortBy, sortOrder string, user models.User) ([]models.BankAccount, int, error) {
	if m.errGetAll != nil {
		return nil, 0, m.errGetAll
	}
	return m.bankAccountList, m.totalCount, nil
}

// GetUserById simulates retrieving users by their IDs.
func (m *MockBankAccountRepository) GetUserById(ids []int, users *[]models.User) error {
	if m.errGetUser != nil {
		return m.errGetUser
	}
	return nil
}

func (m *MockBankAccountRepository) DeleteBankAccount(id uint) error {
	if m.errorDelete != nil {
		return m.errorDelete
	}

	return nil
}

func (m *MockBankAccountRepository) FindUsingBankAccountOnBilling(bankAccountId int) bool {
	return m.isAccountExist
}

func TestCreateBankAccount(t *testing.T) {
	tests := []struct {
		name           string
		req            *request.BankAccountCreateRequest
		userID         uint
		mockRepo       MockBankAccountRepository
		expectedErr    string
		expectedResult *models.BankAccount
	}{
		{
			name: "Success - Create Bank Account",
			req: &request.BankAccountCreateRequest{
				SchoolID:      1,
				BankName:      "Bank A",
				AccountName:   "John Doe",
				AccountNumber: "1234567890",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				exists: false,
			},
			expectedResult: &models.BankAccount{SchoolID: 1, BankName: "Bank A", AccountName: "John Doe", AccountNumber: "1234567890", Master: models.Master{CreatedBy: 1}},
		},
		{
			name: "Error - Account Number Exists",
			req: &request.BankAccountCreateRequest{
				AccountNumber: "1234567890",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				exists: true,
			},
			expectedErr: "account number already exists",
		},
		{
			name: "Error - Failed to Check Existence",
			req: &request.BankAccountCreateRequest{
				AccountNumber: "0987654321",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				errExists: fmt.Errorf("database error"),
			},
			expectedErr: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := BankAccountService{bankAccountRepository: &tc.mockRepo}
			result, err := service.CreateBankAccount(tc.req, tc.userID)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	tests := []struct {
		name           string
		bankID         uint
		req            *request.BankAccountCreateRequest
		userID         int
		mockRepo       MockBankAccountRepository
		expectedErr    string
		expectedResult *models.BankAccount
	}{
		{
			name:   "Success - Update Bank Account",
			bankID: 1,
			req: &request.BankAccountCreateRequest{
				SchoolID:      2,
				BankName:      "New Bank",
				AccountName:   "Updated Name",
				AccountNumber: "9876543210",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				bankAccountByID: &models.BankAccount{SchoolID: 1, BankName: "Old Bank", AccountName: "Old Name", AccountNumber: "1234567890"},
			},
			expectedResult: &models.BankAccount{
				SchoolID:      2,
				BankName:      "New Bank",
				AccountName:   "Updated Name",
				AccountNumber: "9876543210",
				Master:        models.Master{UpdatedBy: 1},
			},
		},
		{
			name:   "Error - Account Number Exists",
			bankID: 1,
			req: &request.BankAccountCreateRequest{
				AccountNumber: "9876543210",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				exists: true,
			},
			expectedErr: "account number already exists",
		},
		{
			name:   "Error - Bank Account Not Found",
			bankID: 2,
			req: &request.BankAccountCreateRequest{
				AccountNumber: "1122334455",
			},
			userID: 1,
			mockRepo: MockBankAccountRepository{
				errGetByID: fmt.Errorf("bank account not found"),
			},
			expectedErr: "bank account not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := BankAccountService{bankAccountRepository: &tc.mockRepo}
			result, err := service.UpdateBankAccount(tc.bankID, tc.req, tc.userID)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

// func TestGetAllBankAccounts(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		page          int
// 		limit         int
// 		search        string
// 		mockRepo      MockBankAccountRepository
// 		expectedErr   string
// 		expectedResult response.BankAccountListResponse
// 		userId          int
// 	}{
// 		{
// 			name: "Success - Get All Bank Accounts",
// 			page: 1,
// 			limit: 10,
// 			mockRepo: MockBankAccountRepository{
// 				bankAccountList: []models.BankAccount{
// 					{BankName: "Bank A", AccountName: "Account A", AccountNumber: "1234567890", AccountOwner: "test", Master: models.Master{CreatedBy: 1}, School: &models.School{SchoolName: "School A"}},
// 					{BankName: "Bank B", AccountName: "Account B", AccountNumber: "0987654321", AccountOwner: "test", Master: models.Master{CreatedBy: 2}, School: &models.School{SchoolName: "School B"}},
// 				},
// 				totalCount: 2,
// 			},
// 			expectedResult: response.BankAccountListResponse{
// 				Limit:      10,
// 				Page:       1,
// 				TotalData:  2,
// 				TotalPage:  1,
// 				Data: []response.BankAccountData{
// 					{BankName: "Bank A", AccountName: "Account A", AccountNumber: "1234567890", AccountOwner: "test", CreatedBy: "", School: response.SchoolData{Name: "School A"}, PlaceHolder: "Bank A-Account A-1234567890"},
// 					{BankName: "Bank B", AccountName: "Account B", AccountNumber: "0987654321", AccountOwner: "test", CreatedBy: "", School: response.SchoolData{Name: "School B"}, PlaceHolder: "Bank B-Account B-0987654321"},
// 				},
// 			},
// 		},
// 		{
// 			name: "Error - Failed to Get All Bank Accounts",
// 			page: 0,
// 			limit: 0,
// 			mockRepo: MockBankAccountRepository{
// 				errGetAll: fmt.Errorf("database error"),
// 			},
// 			expectedErr: "database error",
// 		},
// 		{
// 			name: "Success - Search Functionality",
// 			page: 1,
// 			limit: 10,
// 			search: "Bank A",
// 			mockRepo: MockBankAccountRepository{
// 				bankAccountList: []models.BankAccount{
// 					{BankName: "Bank A", AccountName: "Account A", AccountNumber: "1234567890", AccountOwner: "test", Master: models.Master{CreatedBy: 1}, School: &models.School{SchoolName: "School A"}},
// 				},
// 				totalCount: 1,
// 			},
// 			expectedResult: response.BankAccountListResponse{
// 				Limit:      10,
// 				Page:       1,
// 				TotalData:  1,
// 				TotalPage:  1,
// 				Data: []response.BankAccountData{
// 					{BankName: "Bank A", AccountName: "Account A", AccountNumber: "1234567890", AccountOwner: "test", CreatedBy: "", School: response.SchoolData{Name: "School A"}, PlaceHolder: "Bank A-Account A-1234567890"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			service := BankAccountService{bankAccountRepository: &tc.mockRepo}
// 			result, err := service.GetAllBankAccounts(tc.page, tc.limit, tc.search, "", "", tc.userId)

// 			if tc.expectedErr != "" {
// 				assert.ErrorContains(t, err, tc.expectedErr)
// 				assert.Equal(t, response.BankAccountListResponse{}, result)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tc.expectedResult, result)
// 			}
// 		})
// 	}
// }

func TestGetBankAccountDetails(t *testing.T) {
	tests := []struct {
		name           string
		id             uint
		mockRepo       MockBankAccountRepository
		expectedErr    string
		expectedResult *response.BankAccountData
	}{
		{
			name: "Success - Get Bank Account Details",
			id:   1,
			mockRepo: MockBankAccountRepository{
				bankAccountByID: &models.BankAccount{
					BankName:      "Bank A",
					AccountName:   "Account A",
					AccountNumber: "1234567890",
					AccountOwner:  "Owner A",
					Master:        models.Master{CreatedBy: 1},
					School: &models.School{
						SchoolName: "School A",
					},
				},
			},
			expectedResult: &response.BankAccountData{
				BankName:      "Bank A",
				AccountName:   "Account A",
				AccountNumber: "1234567890",
				AccountOwner:  "Owner A",
				CreatedBy:     "", // Assuming the mocked user returns "testuser"
				School: response.SchoolData{
					Name: "School A",
				},
				PlaceHolder: "Bank A-Account A-1234567890",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := BankAccountService{bankAccountRepository: &tc.mockRepo}

			result, err := service.GetBankAccountDetails(tc.id)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDeleteBankAccount(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockRepo    MockBankAccountRepository
		expectedErr string
	}{
		{
			name: "Success - Delete Bank Account",
			id:   1,
			mockRepo: MockBankAccountRepository{
				bankAccountByID: &models.BankAccount{BankName: "Bank A"}, // Simulate found bank account
				errGetByID:      nil,                                     // No error for getting by ID
				errorDelete:     nil,                                     // No error for deletion
			},
		},
		{
			name: "Error - Bank Account Not Found",
			id:   2,
			mockRepo: MockBankAccountRepository{
				errorDelete: fmt.Errorf("bank account not found"), // Simulate not found error
			},
			expectedErr: "bank account not found",
		},
		{
			name: "Error - Internal Error on Deletion",
			id:   3,
			mockRepo: MockBankAccountRepository{
				bankAccountByID: &models.BankAccount{BankName: "Bank A"}, // Simulate found bank account
				errorDelete:     fmt.Errorf("database error"),            // Simulate a database error
			},
			expectedErr: "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := BankAccountService{bankAccountRepository: &tc.mockRepo}
			err := service.DeleteBankAccount(tc.id)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
