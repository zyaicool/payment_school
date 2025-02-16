package repositories

import (
	"strings"
	"testing"
	"time"

	"schoolPayment/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetAllSchoolList(t *testing.T) {
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

	repo := NewSchoolRepository(gormDB)

	// Define mock data
	searchTerm := "Test School"
	page, limit := 1, 10

	// Mock the count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "schools" LEFT JOIN users ON users.id = schools.created_by WHERE schools.deleted_at IS NULL AND lower\(schools\.school_name\) LIKE`).
		WithArgs("%" + strings.ToLower(searchTerm) + "%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock the select query with the updated order clause and both arguments
	rows := sqlmock.NewRows([]string{"id", "school_name", "created_by_username"}).
		AddRow(1, "Test School", "system")

	mock.ExpectQuery(`SELECT schools\.\*, case when users.username is null then 'system' else users.username end AS created_by_username FROM "schools" LEFT JOIN users ON users.id = schools.created_by WHERE schools.deleted_at IS NULL AND lower\(schools\.school_name\) LIKE`).
		WithArgs("%"+strings.ToLower(searchTerm)+"%", limit).
		WillReturnRows(rows)

	// Call GetAllSchoolList
	schoolList, count, err := repo.GetAllSchoolList(page, limit, searchTerm, "", "")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, schoolList, 1)
	assert.Equal(t, "Test School", schoolList[0].SchoolName)
	assert.Equal(t, "system", schoolList[0].CreatedByUsername)
}

func TestGetSchoolByID(t *testing.T) {
	// Membuat sqlmock database dan koneksi GORM
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Gunakan konfigurasi yang benar untuk GORM dan sqlmock
	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// Pastikan objek `gormDB` tidak `nil`
	if gormDB == nil {
		t.Fatalf("gormDB is nil, failed to initialize")
	}

	repo := NewSchoolRepository(gormDB)

	// Mock query
	rows := sqlmock.NewRows([]string{"id", "school_name"}).AddRow(1, "Test School")
	// Mock query yang sesuai dengan query yang dijalankan GORM
	mock.ExpectQuery(`SELECT \* FROM "schools" WHERE id = \$1 AND deleted_at IS NULL ORDER BY "schools"."id" LIMIT \$2`).
		WithArgs(1, 1). // Tambahkan dua argumen: `id` dan `LIMIT 1`
		WillReturnRows(rows)

	school, err := repo.GetSchoolByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), school.ID)
	assert.Equal(t, "Test School", school.SchoolName)
}

func TestUpdateSchool(t *testing.T) {
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

	repo := NewSchoolRepository(gormDB)

	// Define the school to update
	school := &models.School{
		ID:         1,
		SchoolName: "Updated School Name",
	}

	// Mock the SQL query for update
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "schools" SET "created_at"=\$1,"created_by"=\$2,"updated_at"=\$3,"updated_by"=\$4,"deleted_at"=\$5,"deleted_by"=\$6,"npsn"=\$7,"school_code"=\$8,"school_name"=\$9,"school_province"=\$10,"school_city"=\$11,"school_phone"=\$12,"school_address"=\$13,"school_mail"=\$14,"school_fax"=\$15,"school_logo"=\$16,"school_letterhead"=\$17,"school_grade_id"=\$18 WHERE "id" = \$19`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // created_by
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // updated_by
			sqlmock.AnyArg(), // deleted_at
			sqlmock.AnyArg(), // deleted_by
			school.Npsn,
			school.SchoolCode,
			school.SchoolName,
			school.SchoolProvince,
			school.SchoolCity,
			school.SchoolPhone,
			school.SchoolAddress,
			school.SchoolMail,
			school.SchoolFax,
			school.SchoolLogo,
			school.SchoolLetterhead,
			school.SchoolGradeID,
			school.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the UpdateSchool method
	updatedSchool, err := repo.UpdateSchool(school)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Updated School Name", updatedSchool.SchoolName)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateSchool tests the CreateSchool function.
func TestCreateSchool(t *testing.T) {
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
	mock.ExpectQuery(`INSERT INTO "schools"`).
		WithArgs(
			createdAt,          // created_at
			0,                  // created_by
			updatedAt,          // updated_at
			0,                  // updated_by
			nil,                // deleted_at
			nil,                // deleted_by
			0,                  // npsn
			"",                 // school_code
			"Test School",      // school_name
			"",                 // school_province
			"",                 // school_city
			"",                 // school_phone
			"123 Test Address", // school_address
			"",                 // school_mail
			"",                 // school_fax
			"",                 // school_logo
			"",                 // school_letterhead
			0,                  // school_grade_id
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Initialize repository with mocked DB
	repo := &schoolRepository{
		db: gormDB,
	}

	// Create a test school object
	school := &models.School{
		Master: models.Master{
			CreatedAt: createdAt,
			CreatedBy: 0, // Assuming 0 is the default
			UpdatedAt: updatedAt,
			UpdatedBy: 0,
			DeletedAt: nil, // Assuming nil is the default
			DeletedBy: nil,
		},
		Npsn:             0,  // Assuming 0 is the default
		SchoolCode:       "", // Default empty string
		SchoolName:       "Test School",
		SchoolAddress:    "123 Test Address",
		SchoolProvince:   "",
		SchoolCity:       "",
		SchoolPhone:      "",
		SchoolMail:       "",
		SchoolFax:        "",
		SchoolLogo:       "",
		SchoolLetterhead: "",
		SchoolGradeID:    0,
	}

	// Call CreateSchool
	result, err := repo.CreateSchool(school)
	assert.NoError(t, err)
	assert.Equal(t, school, result)

	// Ensure all expectations are met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
