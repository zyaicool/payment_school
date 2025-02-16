package repositories

import (
	"errors"
	"regexp"
	"testing"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StudentRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *gorm.DB
	repo StudentRepositoryInteface
}

func (s *StudentRepositoryTestSuite) SetupTest() {
	var err error
	var db *gorm.DB

	// Create mock db
	mockDB, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	database.DB = db
	s.mock = mock
	s.db = db
	s.repo = NewStudentRepository(db)
}

func TestStudentRepositorySuite(t *testing.T) {
	suite.Run(t, new(StudentRepositoryTestSuite))
}

func (s *StudentRepositoryTestSuite) TestGetAllStudent() {
	// Test cases
	testCases := []struct {
		name          string
		page          int
		limit         int
		search        string
		user          models.User
		status        string
		gradeID       int
		yearId        int
		schoolID      int
		searchNis     string
		classID       int
		sortBy        string
		sortOrder     string
		studentId     int
		isActive      *bool
		expectedCount int64
		expectError   bool
	}{
		{
			name:          "Success with basic filters",
			page:          1,
			limit:         10,
			search:        "",
			user:          models.User{RoleID: 1},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "Success with search",
			page:          1,
			limit:         10,
			search:        "John",
			user:          models.User{RoleID: 1},
			expectedCount: 1,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// Mock count query
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(tc.expectedCount)
			s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
				WillReturnRows(countRows)

			// Mock data query
			rows := sqlmock.NewRows([]string{
				"id", "full_name", "nis", "school_grade", "school_class",
				"school_year_name", "placeholder",
			}).
				AddRow(1, "John Doe", "12345", "Grade 10", "Class A", "2023/2024", "12345 - John Doe, Class A, Grade 10").
				AddRow(2, "Jane Doe", "12346", "Grade 10", "Class A", "2023/2024", "12346 - Jane Doe, Class A, Grade 10")

			s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
				WillReturnRows(rows)

			students, _, total, err := GetAllStudent(
				tc.page, tc.limit, tc.search, tc.user, tc.status,
				tc.gradeID, tc.yearId, tc.schoolID, tc.searchNis,
				tc.classID, tc.sortBy, tc.sortOrder, tc.studentId, tc.isActive,
			)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, total)
				assert.NotNil(t, students)
			}
		})
	}
}

func (s *StudentRepositoryTestSuite) TestGetStudentByID() {
	testCases := []struct {
		name        string
		id          uint
		user        models.User
		expectError bool
	}{
		{
			name:        "Success get student",
			id:          1,
			user:        models.User{RoleID: 1},
			expectError: false,
		},
		{
			name:        "Student not found",
			id:          999,
			user:        models.User{RoleID: 1},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			if !tc.expectError {
				rows := sqlmock.NewRows([]string{"id", "full_name", "nis"}).
					AddRow(1, "John Doe", "12345")
				s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
					WithArgs(tc.id, 1).
					WillReturnRows(rows)
			} else {
				s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
					WithArgs(tc.id, 1).
					WillReturnError(errors.New("record not found"))
			}

			student, err := s.repo.GetStudentByID(tc.id, tc.user)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.id, student.ID)
			}
		})
	}
}

func (s *StudentRepositoryTestSuite) TestCreateStudent() {
	student := &models.Student{
		FullName: "John Doe",
		Nis:      "12345",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "students"`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // created_by
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // updated_by
			nil,              // deleted_at
			nil,              // deleted_by
			"",               // nisn
			"",               // registration_number
			student.Nis,      // nis
			"",               // nik
			student.FullName, // full_name
			"",               // gender
			"",               // religion
			"",               // citizenship
			"",               // birth_place
			nil,              // birth_date
			"",               // address
			"",               // school_grade
			0,                // school_grade_id
			"",               // school_class
			0,                // school_class_id
			"",               // no_handphone
			"",               // height
			"",               // weight
			"",               // medical_history
			0,                // distance_to_school
			"",               // sibling
			"",               // nick_name
			"",               // email
			"",               // entry_year
			"",               // status
			"",               // image
			0,                // school_year_id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	result, err := s.repo.CreateStudent(student)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), student.FullName, result.FullName)
}

func (s *StudentRepositoryTestSuite) TestUpdateStudent() {
	student := &models.Student{
		Master: models.Master{
			ID: 1,
		},
		FullName: "John Doe Updated",
		Nis:      "12345",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	result, err := s.repo.UpdateStudent(student)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), student.FullName, result.FullName)
}

func (s *StudentRepositoryTestSuite) TestBulkCreateStudents() {
	students := []models.Student{
		{
			FullName: "John Doe",
			Nis:      "12345",
		},
		{
			FullName: "Jane Doe",
			Nis:      "12346",
		},
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "students"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))
	s.mock.ExpectCommit()

	err := s.repo.BulkCreateStudents(students)
	assert.NoError(s.T(), err)
}

func (s *StudentRepositoryTestSuite) TestGetStudentByNis() {
	testCases := []struct {
		name        string
		nis         string
		user        models.User
		expectError bool
	}{
		{
			name:        "Success get student by NIS",
			nis:         "12345",
			user:        models.User{RoleID: 1},
			expectError: false,
		},
		{
			name:        "Student not found",
			nis:         "99999",
			user:        models.User{RoleID: 1},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			if !tc.expectError {
				rows := sqlmock.NewRows([]string{"id", "full_name", "nis"}).
					AddRow(1, "John Doe", tc.nis)
				s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
					WillReturnRows(rows)
			} else {
				s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
					WillReturnError(errors.New("record not found"))
			}

			student, err := s.repo.GetStudentByNis(tc.nis, tc.user)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.nis, student.Nis)
			}
		})
	}
}
