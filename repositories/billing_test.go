package repositories_test

import (
	"fmt"
	"testing"
	"time"

	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Mock the BillingRepositoryInterface
type MockBillingRepository struct {
	mock.Mock
	repositories.BillingRepositoryInterface 
}

func (m *MockBillingRepository) GetDetailBillingsByBillingID(billingID uint) ([]response.DetailBilling, error) {
	args := m.Called(billingID)
	// Mengembalikan nilai yang tepat meskipun nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]response.DetailBilling), args.Error(1)
}

// Now, we can write the unit test

func TestGetDetailBillingsByBillingID(t *testing.T) {
	// Initialize the mock
	mockRepo := new(MockBillingRepository)

	// Define expected results
	expectedDetailBillings := []response.DetailBilling{
		{
			ID:                1,
			DetailBillingName: "Tuition Fee",
			DueDate:           time.Now(),
			Amount:            1000,
		},
		{
			ID:                2,
			DetailBillingName: "Lab Fee",
			DueDate:           time.Now(),
			Amount:            500,
		},
	}

	// Set up expectations for the mock
	mockRepo.On("GetDetailBillingsByBillingID", uint(1)).Return(expectedDetailBillings, nil)

	// Call the method
	result, err := mockRepo.GetDetailBillingsByBillingID(1)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedDetailBillings, result)

	// Ensure expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetDetailBillingsByBillingIDError(t *testing.T) {
	// Initialize the mock
	mockRepo := new(MockBillingRepository)

	// Set up mock to simulate an error
	mockRepo.On("GetDetailBillingsByBillingID", uint(1)).Return(nil, fmt.Errorf("database error"))

	// Call the method
	result, err := mockRepo.GetDetailBillingsByBillingID(1)

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)

	// Ensure expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateBillingAccount(t *testing.T) {
	// Test cases
	testCases := []struct {
		name       string
		input      *models.Billing
		mockResult func(mock sqlmock.Sqlmock)
		expectErr  bool
	}{
		{
			name: "Successful CreateBilling",
			input: &models.Billing{
				Master: models.Master{
					CreatedAt: mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"),
					CreatedBy: 0,
					UpdatedAt: mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"),
					UpdatedBy: 0,
				},
				BillingNumber:  "BN001",
				BillingName:    "Test Billing",
				BillingType:    "Type1",
				SchoolGradeID:  1,
				SchoolYearId:   1,
				BillingAmount:  1000,
				Description:    "Test Description",
				BillingCode:    "BC001",
				SchoolClassIds: "[1,2]",
				BankAccountId:  1,
				IsDonation:     false,
			},
			mockResult: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "billings" \("created_at","created_by","updated_at","updated_by","deleted_at","deleted_by","billing_number","billing_name","billing_type","school_grade_id","school_year_id","billing_amount","description","billing_code","school_class_ids","bank_account_id","is_donation"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14,\$15,\$16,\$17\) RETURNING "id"`).
					WithArgs(
						mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"), // created_at
						0,                 // created_by
						mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"), // updated_at
						0,                 // updated_by
						nil,               // deleted_at
						nil,               // deleted_by
						"BN001",           // billing_number
						"Test Billing",    // billing_name
						"Type1",           // billing_type
						1,                 // school_grade_id
						1,                 // school_year_id
						1000,              // billing_amount
						"Test Description",// description
						"BC001",           // billing_code
						"[1,2]",           // school_class_ids
						1,                 // bank_account_id
						false,             // is_donation
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Fail CreateBilling - DB Error",
			input: &models.Billing{
				Master: models.Master{
					CreatedAt: mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"),
					CreatedBy: 0,
					UpdatedAt: mustParseTime("2006-01-02 15:04:05", "2023-01-01 00:00:00"),
					UpdatedBy: 0,
				},
				BillingNumber:  "BN002",
				BillingName:    "Error Billing",
				BillingType:    "Type2",
				SchoolGradeID:  2,
				SchoolYearId:   2,
				BillingAmount:  2000,
				Description:    "Error Description",
				BillingCode:    "BC002",
				SchoolClassIds: "[3,4]",
				BankAccountId:  2,
				IsDonation:     true,
			},
			mockResult: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "billings"`).WillReturnError(fmt.Errorf("DB error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize sqlmock and gorm DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			assert.NoError(t, err)

			// Setup mock behavior
			tc.mockResult(mock)

			// Initialize repository with mocked DB
			repo := repositories.NewBillingRepository(gormDB)

			// Call CreateBilling
			result, err := repo.CreateBilling(tc.input)

			// Validate results
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, result)
			}

			// Ensure all expectations are met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

// Helper function to parse time and handle errors
func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestGetAllBilling_Success(t *testing.T) {
	// Initialize sqlmock and gorm DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create mock repository
	repo := repositories.NewBillingRepository(gormDB)

	// Input parameters
	page := 1
	limit := 10
	search := "Test"
	billingType := "Type1"
	paymentType := "Online"
	schoolGrade := "1"
	sort := "ASC"
	sortBy := "created_at"
	sortOrder := "ASC"
	bankAccountId := 1
	isDonation := false
	user := models.User{
		UserSchool: &models.UserSchool{
			SchoolID: 1,
		},
	}

	// Mock count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "billings"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

	// Mock main query
	mock.ExpectQuery(`SELECT billings\.\*, CASE WHEN billings\.is_donation = true THEN 'admin' ELSE users\.username END as create_by_username FROM "billings"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "billing_name", "billing_type", "school_grade_id", "school_year_id", "is_donation", "bank_account_id",
		}).AddRow(1, "Test Billing", "Type1", 1, 2023, false, 1))

	// Mock Preload for BankAccount
	mock.ExpectQuery(`SELECT \* FROM "bank_accounts" WHERE "bank_accounts"\."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "account_number",
		}).AddRow(1, "123456789"))

	// Mock Preload for SchoolGrade
	mock.ExpectQuery(`SELECT \* FROM "school_grades" WHERE "school_grades"\."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Grade 1"))

	// Call function
	billings, totalPages, total, err := repo.GetAllBilling(
		page,
		limit,
		search,
		billingType,
		paymentType,
		schoolGrade,
		sort,
		sortBy,
		sortOrder,
		bankAccountId,
		&isDonation,
		user,
	)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, billings)
	assert.Equal(t, int64(100), total)      // Total records
	assert.Equal(t, 10, totalPages)        // Total pages
	assert.Len(t, billings, 1)             // One record returned
	assert.Equal(t, "Test Billing", billings[0].BillingName)
	assert.NotNil(t, billings[0].SchoolGrade)
	assert.NotNil(t, billings[0].BankAccount)
	assert.Equal(t, "123456789", billings[0].BankAccount.AccountNumber)

	// Ensure all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
