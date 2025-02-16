package repositories_test

import (
	"fmt"
	"testing"

	"schoolPayment/dtos/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the BillingHistoryRepository
type MockBillingHistoryRepository struct {
	mock.Mock
}

func (m *MockBillingHistoryRepository) GetAllBillingHistory(page int, limit int, search string, studentId int, roleID int, schoolYearId int, paymentTypeId int, schoolID int, paymentStatusCode string, sortBy string, sortOrder string, userID int, userLoginID int) ([]response.DataListBillingHistory, int, int64, error) {
	args := m.Called(page, limit, search, studentId, roleID, schoolYearId, paymentTypeId, schoolID, paymentStatusCode, sortBy, sortOrder, userID, userLoginID)
	return args.Get(0).([]response.DataListBillingHistory), args.Int(1), args.Get(2).(int64), args.Error(3)
}

func TestGetAllBillingHistory(t *testing.T) {

    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 1
    totalData := int64(1)

    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).Return(expectedBillingHistory, totalPages, totalData, nil)
    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}

func TestGetAllBillingHistoryWithError(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Simulate an error being returned
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return([]response.DataListBillingHistory{}, 0, int64(0), fmt.Errorf("database error"))

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.NotNil(t, err)
    assert.Equal(t, "database error", err.Error())
    assert.Empty(t, result)
    assert.Equal(t, 0, totalPages)
    assert.Equal(t, int64(0), totalData)

    mockRepo.AssertExpectations(t)
}

//Uji skenario pencarian (search) berdasarkan nis atau full_name dalam parameter query.
func TestGetAllBillingHistoryWithSearch(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Test data expected
    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 1
    totalData := int64(1)

    // Simulate a successful query with a search term
    mockRepo.On("GetAllBillingHistory", 1, 10, "john", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, totalPages, totalData, nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "john", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}

//Uji skenario pagination berdasarkan parameter page dan limit, dan periksa jika perhitungan totalPages benar.
func TestGetAllBillingHistoryWithPagination(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 2
    totalData := int64(15)

    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, totalPages, totalData, nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}

//Uji skenario dengan filter berdasarkan roleID dan verifikasi apakah query tambahan ditambahkan sesuai dengan role.
func TestGetAllBillingHistoryWithRoleFilter(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 1
    totalData := int64(1)

    // Simulate a query with roleID = 2 (user role)
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, totalPages, totalData, nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}
//Uji kondisi ketika filter paymentStatusCode digunakan dalam query.
func TestGetAllBillingHistoryWithPaymentStatus(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 1
    totalData := int64(1)

    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "PS02", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, totalPages, totalData, nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "PS02", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}

//Uji skenario dengan parameter limit = 0, yang akan mengabaikan pagination dan mengembalikan semua data dalam satu respons.
func TestGetAllBillingHistoryWithLimitZero(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    totalPages := 1
    totalData := int64(1)

    mockRepo.On("GetAllBillingHistory", 1, 0, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, totalPages, totalData, nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 0, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, totalPages, totalPages)
    assert.Equal(t, totalData, totalData)

    mockRepo.AssertExpectations(t)
}

//Uji skenario error pada database untuk memastikan error handling bekerja dengan baik.

func TestGetAllBillingHistoryWithDatabaseError(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Simulate a database error scenario by returning a valid empty slice and the error
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return([]response.DataListBillingHistory{}, 0, int64(0), fmt.Errorf("database error"))

    // Call the method
    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    // Assertions
    assert.NotNil(t, err)
    assert.Equal(t, "database error", err.Error())
    assert.Empty(t, result)
    assert.Equal(t, 0, totalPages)
    assert.Equal(t, int64(0), totalData)

    // Ensure expectations were met
    mockRepo.AssertExpectations(t)
}

// Test untuk skenario dengan roleID yang berbeda
func TestGetAllBillingHistoryWithRoleID(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Simulasi roleID yang berbeda
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return([]response.DataListBillingHistory{}, 0, int64(0), nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Empty(t, result)
    assert.Equal(t, 0, totalPages)
    assert.Equal(t, int64(0), totalData)
    mockRepo.AssertExpectations(t)
}

// Test untuk pagination dengan limit
func TestGetAllBillingHistoryWithPaginationAndLimit(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    expectedBillingHistory := []response.DataListBillingHistory{
        {
            ID:                1,
            InvoiceNumber:     "INV123",
            StudentName:       "John Doe",
            PaymentDate:       nil,
            PaymentMethod:     "Credit Card",
            Username:          "admin",
            TotalAmount:       1000,
            TransactionStatus: "Paid",
            OrderID:           "ORD123",
            Token:             "token123",
            RedirectUrl:       "http://example.com",
        },
    }

    // Test limit > 0 and pagination
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return(expectedBillingHistory, 1, int64(1), nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, expectedBillingHistory, result)
    assert.Equal(t, 1, totalPages)
    assert.Equal(t, int64(1), totalData)
    mockRepo.AssertExpectations(t)
}

func TestGetAllBillingHistoryWithRoleAndPaymentStatus(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Simulasi dengan Role ID 1 dan Payment Status "Paid"
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 1, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return([]response.DataListBillingHistory{}, 0, int64(0), nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 1, 2, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Empty(t, result)
    assert.Equal(t, 0, totalPages)
    assert.Equal(t, int64(0), totalData)
    mockRepo.AssertExpectations(t)
}

// func TestGetAllBillingHistoryWithErrorHandling(t *testing.T) {
//     mockRepo := new(MockBillingHistoryRepository)

//     // Simulasi database error
//     mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
//         Return(nil, 0, 0, fmt.Errorf("database error"))

//     result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

//     assert.NotNil(t, err)
//     assert.Equal(t, "database error", err.Error())
//     assert.Empty(t, result)
//     assert.Equal(t, 0, totalPages)
//     assert.Equal(t, int64(0), totalData)
//     mockRepo.AssertExpectations(t)
// }

func TestGetAllBillingHistoryWithPaginationEdgeCase(t *testing.T) {
    mockRepo := new(MockBillingHistoryRepository)

    // Simulasi pagination edge case
    mockRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0).
        Return([]response.DataListBillingHistory{
            {ID: 1, InvoiceNumber: "INV123", TotalAmount: 1000},
        }, 1, int64(1), nil)

    result, totalPages, totalData, err := mockRepo.GetAllBillingHistory(1, 10, "", 0, 1, 2024, 1, 1, "", "created_at", "asc", 0, 0)

    assert.Nil(t, err)
    assert.Equal(t, 1, totalPages)
    assert.Equal(t, int64(1), totalData)
    assert.Len(t, result, 1)
    mockRepo.AssertExpectations(t)
}