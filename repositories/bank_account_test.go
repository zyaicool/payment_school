package repositories

import (
	"fmt"
	"schoolPayment/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUpdateBankAccount(t *testing.T) {
	// Set up sqlmock database and GORM connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Configure GORM with sqlmock
	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// Ensure `gormDB` is not nil
	if gormDB == nil {
		t.Fatalf("gormDB is nil, failed to initialize")
	}

	repo := NewBankAccountRepository(gormDB)

	// Define the bank account to update
	bankAccount := &models.BankAccount{
		Master: models.Master{
			ID: 1,
		},
		BankName: "Bank Name",
	}

	// Mock the SQL query for update
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "bank_accounts" SET "created_at"=\$1,"created_by"=\$2,"updated_at"=\$3,"updated_by"=\$4,"deleted_at"=\$5,"deleted_by"=\$6,"school_id"=\$7,"bank_name"=\$8,"account_name"=\$9,"account_number"=\$10,"account_owner"=\$11 WHERE "id" = \$12`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // created_by
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // updated_by
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // deleted_by
			bankAccount.SchoolID,
			bankAccount.BankName,
			bankAccount.AccountName,
			bankAccount.AccountNumber,
			bankAccount.AccountOwner,
			bankAccount.Master.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the UpdateBankAccount method
	updatedBankAccount, err := repo.UpdateBankAccount(bankAccount)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Bank Name", updatedBankAccount.BankName) // Make sure to check the expected result here

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateBankAccount(t *testing.T) {
	// Initialize sqlmock and gorm DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Parse dates for CreatedAt and UpdatedAt
	createdAt, err := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	assert.NoError(t, err)
	updatedAt := createdAt // Use the same time for simplicity

	// Prepare mock expectation
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "bank_accounts"`).
		WithArgs(
			createdAt,          // created_at
			0,                  // created_by
			updatedAt,          // updated_at
			0,                  // updated_by
			nil,                // deleted_at
			nil,                // deleted_by
			1,
			"Bank Account",
			"",
			"",
			"",
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Initialize repository with mocked DB
	repo := &bankAccountRepository{
		db: gormDB,
	}

	// Create a test school object
	bankAccount := &models.BankAccount{
		Master: models.Master{
			CreatedAt: createdAt,
			CreatedBy: 0, // Assuming 0 is the default
			UpdatedAt: updatedAt,
			UpdatedBy: 0,
			DeletedAt: nil, // Assuming nil is the default
			DeletedBy: nil,
		},
		SchoolID: 1,
		BankName: "Bank Account",
		AccountName: "",
		AccountNumber: "",
		AccountOwner: "",
	}

	// Call CreateSchool
	result, err := repo.CreateBankAccount(bankAccount)
	assert.NoError(t, err)
	assert.Equal(t, bankAccount, result)

	// Ensure all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteBankAccount(t *testing.T) {
	// Initialize sqlmock and gorm DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		name            string
		id              uint
		rowsAffected    int64
		expectedError   error
		prepareMock     func(id uint, rowsAffected int64)
	}{
		{
			name:         "successful deletion",
			id:           1,
			rowsAffected: 1,
			expectedError: nil,
			prepareMock: func(id uint, rowsAffected int64) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "bank_accounts" WHERE "bank_accounts"."id" = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, rowsAffected))
				mock.ExpectCommit()
			},
		},
		{
			name:         "no rows affected",
			id:           2,
			rowsAffected: 0,
			expectedError: nil, // Assuming no error if no rows are deleted
			prepareMock: func(id uint, rowsAffected int64) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "bank_accounts" WHERE "bank_accounts"."id" = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, rowsAffected))
				mock.ExpectCommit()
			},
		},
		{
			name:         "SQL error",
			id:           3,
			rowsAffected: 0,
			expectedError: fmt.Errorf("mock SQL error"),
			prepareMock: func(id uint, rowsAffected int64) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "bank_accounts" WHERE "bank_accounts"."id" = \$1`).
					WithArgs(id).
					WillReturnError(fmt.Errorf("mock SQL error"))
				mock.ExpectRollback()
			},
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare mock expectations
			tc.prepareMock(tc.id, tc.rowsAffected)

			// Initialize repository
			repo := &bankAccountRepository{
				db: gormDB,
			}

			// Call DeleteBankAccount
			err := repo.DeleteBankAccount(tc.id)

			// Assertions
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations are met
			mockErr := mock.ExpectationsWereMet()
			assert.NoError(t, mockErr)
		})
	}
}


