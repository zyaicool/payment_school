package repositories

import (
	"regexp"
	database "schoolPayment/configs"
	"schoolPayment/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetGetSchoolYearByIDRepository(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	database.DB = gormDB

	t.Run("Success", func(t *testing.T) {
		id := 28

		// Mengonversi string tanggal ke *time.Time
		startDate, err := time.Parse("2006-01-02 15:04:05.000 -0700", "2026-11-04 07:00:00.000 +0700")
		if err != nil {
			t.Fatalf("failed to parse start date: %v", err)
		}
		endDate, err := time.Parse("2006-01-02 15:04:05.000 -0700", "2027-11-29 07:00:00.000 +0700")
		if err != nil {
			t.Fatalf("failed to parse end date: %v", err)
		}

		// Menyiapkan mock data dengan nilai yang telah dikonversi menjadi *time.Time
		rows := sqlmock.NewRows([]string{"id", "school_year_name", "start_date", "end_date"}).
			AddRow(28, "2026/2027", startDate, endDate)

		// Perbaiki ekspresi reguler untuk mencocokkan query yang sebenarnya
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "school_years" WHERE id = $1 AND deleted_at IS NULL ORDER BY "school_years"."id" LIMIT $2`)).
			WithArgs(id, 1). // Karena LIMIT diquerynya adalah 1
			WillReturnRows(rows)

		repo := NewSchoolYearRepository(gormDB)
		result, err := repo.GetSchoolYearByID(uint(id))

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, "2026/2027", result.SchoolYearName)
		assert.NotNil(t, result.StartDate)
		assert.NotNil(t, result.EndDate)
	})
}

func TestCreateSchoolYearRespository(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	database.DB = gormDB

	t.Run("Success create year", func(t *testing.T) {
		// Buat repository dengan menyertakan gormDB
		repo := NewSchoolYearRepository(gormDB)

		year := &models.SchoolYear{
			Master: models.Master{
				CreatedBy: 1,
				UpdatedBy: 1,
			},
			SchoolYearCode: "SY001",
			SchoolYearName: "2023/2025",
			SchoolId:       1,
			StartDate: func() *time.Time {
				t, _ := time.Parse("2006-01-02", "2023-01-01")
				return &t
			}(),
			EndDate: func() *time.Time {
				t, _ := time.Parse("2006-01-02", "2025-12-31")
				return &t
			}(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "school_years" ("created_at","created_by","updated_at","updated_by","deleted_at","deleted_by","school_year_code","school_year_name","school_id","start_date","end_date") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				year.Master.CreatedBy,
				sqlmock.AnyArg(), // updated_at
				year.Master.UpdatedBy,
				nil, // deleted_at
				nil, // deleted_by
				year.SchoolYearCode,
				year.SchoolYearName,
				year.SchoolId,
				year.StartDate,
				year.EndDate,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		result, err := repo.CreateSchoolYear(year)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, year.SchoolYearName, result.SchoolYearName)
	})
}
