package repositories

import (
	"regexp"
	"testing"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestGetAllPrefixClassRepository(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	database.DB = gormDB

	t.Run("Success with no filters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "prefix_name", "school_id"}).
			AddRow(1, "Prefix1", 1).
			AddRow(2, "Prefix2", 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prefix_classes" WHERE prefix_classes.deleted_at IS NULL ORDER BY created_at DESC`)).
			WillReturnRows(rows)

		result, err := NewPrefixClassRepository().GetAllPrefixClassRepository("", models.User{})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
	})

	t.Run("Success with search filter", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "prefix_name", "school_id"}).
			AddRow(1, "Test Prefix", 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prefix_classes" WHERE prefix_classes.deleted_at IS NULL AND LOWER(prefix_name) like $1 ORDER BY created_at DESC`)).
			WithArgs("%test%").
			WillReturnRows(rows)

		result, err := NewPrefixClassRepository().GetAllPrefixClassRepository("test", models.User{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("Success with school filter", func(t *testing.T) {
		user := models.User{
			UserSchool: &models.UserSchool{
				SchoolID: 1,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "prefix_name", "school_id"}).
			AddRow(1, "School Prefix", 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "prefix_classes" WHERE prefix_classes.deleted_at IS NULL AND prefix_classes.school_id = $1 ORDER BY created_at DESC`)).
			WithArgs(user.UserSchool.SchoolID).
			WillReturnRows(rows)

		result, err := NewPrefixClassRepository().GetAllPrefixClassRepository("", user)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
	})
}

func TestCreatePrefixClassRepository(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	database.DB = gormDB

	t.Run("Success create prefix", func(t *testing.T) {
		prefix := &models.PrefixClass{
			PrefixName: "Test Prefix",
			SchoolID:   1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "prefix_classes" ("created_at","created_by","updated_at","updated_by","deleted_at","deleted_by","prefix_name","school_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil, prefix.PrefixName, prefix.SchoolID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		result, err := NewPrefixClassRepository().CreatePrefixClassRepository(prefix)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, prefix.PrefixName, result.PrefixName)
	})
}

func TestCheckPrefixClassExists(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	database.DB = gormDB

	t.Run("Prefix exists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "prefix_classes" WHERE LOWER(prefix_name) = $1 AND school_id = $2`)).
			WithArgs("existing prefix", uint(1)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		exists, err := NewPrefixClassRepository().CheckPrefixClassExists("Existing Prefix", 1)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Prefix doesn't exist", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "prefix_classes" WHERE LOWER(prefix_name) = $1 AND school_id = $2`)).
			WithArgs("new prefix", uint(1)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		exists, err := NewPrefixClassRepository().CheckPrefixClassExists("New Prefix", 1)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
