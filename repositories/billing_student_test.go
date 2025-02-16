package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock, func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		db.Close()
	}
}

func TestGetInstallmentDetails(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewBillingStudentRepository(gormDB)

	t.Run("Valid student with installments and donations", func(t *testing.T) {
		studentId := 1
		schoolId := uint(1)
	
		// Mock student query
		mock.ExpectQuery(`SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"."id" LIMIT \$2`).
			WithArgs(studentId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(studentId, "aktif"))
	
		// Mock installment query
		mock.ExpectQuery(`SELECT DISTINCT bs.id AS billing_student_id`).
			WithArgs(studentId).
			WillReturnRows(sqlmock.NewRows([]string{
				"billing_student_id", "detail_billing_name", "amount", "due_date", "payment_status", "billing_type", "transaction_status", "updated_at", "created_at",
			}).AddRow(1, "Installment 1", 1000, nil, "belum bayar", "type1", nil, nil, nil))
	
		// Mock donation query
		mock.ExpectQuery(`SELECT DISTINCT b.id as billing_id`).
			WithArgs(studentId, schoolId, schoolId).
			WillReturnRows(sqlmock.NewRows([]string{
				"billing_id", "billing_name", "updated_at", "created_at",
			}).AddRow(1, "Donation 1", nil, nil))
	
		// Call the function
		installments, donations, err := repo.GetInstallmentDetails(studentId, schoolId)
	
		// Assertions
		assert.NoError(t, err)
		assert.Len(t, installments, 1)
		assert.Equal(t, "Installment 1", installments[0].DetailBillingName)
		assert.Len(t, donations, 1)
		assert.Equal(t, "Donation 1", donations[0].BillingName)
	})	

	t.Run("Student with no installments or donations", func(t *testing.T) {
		studentId := 2
		schoolId := uint(2)

		// Mock student query
		mock.ExpectQuery(`SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"."id" LIMIT \$2`).
			WithArgs(studentId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(studentId, "aktif"))

		// Mock installment query
		mock.ExpectQuery(`SELECT DISTINCT bs.id AS billing_student_id`).
			WithArgs(studentId).
			WillReturnRows(sqlmock.NewRows([]string{}))

		// Mock donation query
		mock.ExpectQuery(`SELECT DISTINCT b.id as billing_id`).
			WithArgs(studentId, schoolId, schoolId).
			WillReturnRows(sqlmock.NewRows([]string{}))

		// Call the function
		installments, donations, err := repo.GetInstallmentDetails(studentId, schoolId)

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, installments)
		assert.Empty(t, donations)
	})

	t.Run("Database query error", func(t *testing.T) {
		studentId := 3
		schoolId := uint(3)

		// Mock student query with error
		mock.ExpectQuery(`SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"."id" LIMIT \$2`).
			WithArgs(studentId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Call the function
		installments, donations, err := repo.GetInstallmentDetails(studentId, schoolId)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, installments)
		assert.Nil(t, donations)
	})
}