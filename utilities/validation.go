package utilities

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"schoolPayment/constants"
	"schoolPayment/models"

	"github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
)

// Check if a field is empty
func ValidateFieldNotEmpty(field interface{}, fieldName string) error {
	switch v := field.(type) {
	case string:
		if v == "" {
			return errors.New(fieldName + " cannot be empty")
		}
	case int:
		if v == 0 {
			return errors.New(fieldName + " cannot be zero or empty")
		}
	default:
		return fmt.Errorf("%s has an unsupported type", fieldName)
	}
	return nil
}

// Validate that a field contains only alphanumeric characters (letters and numbers)
func ValidateFieldCombination(field, fieldName string) error {
	// Regex to allow letters, numbers, and spaces
	matched, err := regexp.MatchString("^[a-zA-Z0-9 ]+$", field)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New(fieldName + " can only contain letters and numbers")
	}
	return nil
}

// Validate that a field has a minimal number of characters
func ValidateFieldMinimalCharacters(field, fieldName string, minChars int) error {
	if len(field) < minChars {
		value := fmt.Sprint(minChars)
		return errors.New(fieldName + " must have at least " + string(value) + " characters")
	}
	return nil
}

// Validate that a field has a maximum number of words
func ValidateFieldMaxWords(field, fieldName string, maxWords int) error {
	words := strings.Fields(field)
	if len(words) > maxWords {
		value := fmt.Sprint(maxWords)
		return errors.New(fieldName + " cannot exceed " + value + " words")
	}
	return nil
}

func ValidatePassword(password string) error {
	// Regular expression to check for at least one number
	numberRegex := regexp.MustCompile(`[0-9]`)
	// Regular expression to check for at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	// Regular expression to check for at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	// password min length 8 char
	passwordValidator := validator.New(validator.MinLength(8, nil))

	err := passwordValidator.Validate(password)
	if err != nil {
		return errors.New("Password Minimum length is 8 char.")
	}

	// Check if password has at least one number
	if !numberRegex.MatchString(password) {
		return errors.New("Password must contain at least one number.")
	}

	// Check if password has at least one uppercase letter
	if !uppercaseRegex.MatchString(password) {
		return errors.New("Password must contain at least one uppercase letter.")
	}

	// Check if password has at least one lowercase letter
	if !lowercaseRegex.MatchString(password) {
		return errors.New("Password must contain at least one lowercase letter.")
	}

	return nil // Password is valid
}

func ValidateUsername(username string) error {
	usernameValidator := validator.New(validator.MinLength(5, nil))

	err := usernameValidator.Validate(username)
	if err != nil {
		return errors.New("Username Minimum length is 5 char.")
	}

	return nil
}

func ValidateEmail(email string) error {
	// check format email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)
	validEmail := emailRegex.MatchString(email)
	if !validEmail {
		return errors.New("Email does not match format.")
	}
	return nil
}

func CompareOldPassword(hashedPassword string, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return fmt.Errorf("Password lama Anda tidak sesuai. Mohon periksa kembali.")
	}
	return nil
}

func ComparePassword(oldPassword string, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newPassword))
	if err == nil {
		return fmt.Errorf("Password sudah pernah digunakan")
	}

	return nil
}

func CheckingAccess(user models.User, pageCode string) bool {
	for _, matrix := range user.Role.RoleMatrix {
		if matrix.PageCode == pageCode {
			return true
		}
	}

	return false
}

func ChangeDate(date string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		t := time.Now()
		return &t, err
	}

	return &t, nil
}

func GeneratePassword() (string, error) {
	const (
		lowercase      = "abcdefghijklmnopqrstuvwxyz"
		uppercase      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers        = "0123456789"
		allChars       = lowercase + uppercase + numbers
		passwordLength = 8
	)

	var password strings.Builder

	// Ensure password contains at least one lowercase, one uppercase, and one digit
	mustInclude := []string{
		string(lowercase[randInt(len(lowercase))]),
		string(uppercase[randInt(len(uppercase))]),
		string(numbers[randInt(len(numbers))]),
	}

	// Add the mandatory characters to the password
	for _, ch := range mustInclude {
		password.WriteString(ch)
	}

	// Generate the rest of the password
	for password.Len() < passwordLength {
		randomChar := allChars[randInt(len(allChars))]
		password.WriteString(string(randomChar))
	}

	// Convert to a string
	return shuffleString(password.String()), nil
}

// Helper function to generate a random integer
func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// Helper function to shuffle the password to avoid predictable patterns
func shuffleString(input string) string {
	runes := []rune(input)
	for i := len(runes) - 1; i > 0; i-- {
		j := randInt(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Function to generate a random username of length 5 with letters and numbers
func GenerateRandomUsername(length int) (string, error) {
	username := make([]byte, length)
	for i := range username {
		randomCharIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		username[i] = charset[randomCharIndex.Int64()]
	}
	return string(username), nil
}

func ValidateBillingName(billingName string) error {
	// Validate field not empty
	if err := ValidateFieldNotEmpty(billingName, constants.MessageBillingName); err != nil {
		return err
	}

	// Validate field only contains letters and numbers
	if err := ValidateFieldCombination(billingName, constants.MessageBillingName); err != nil {
		return err
	}

	// Validate field does not exceed 50 words
	if err := ValidateFieldMaxWords(billingName, constants.MessageBillingName, 50); err != nil {
		return err
	}

	return nil
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr) // Adjust the layout as needed
}

// ValidateDateRange checks if the start and end dates are valid and if the end date is after the start date.
func ValidateDateRange(startDateStr, endDateStr string) error {
	startDate, err := ParseDate(startDateStr)
	if err != nil {
		return errors.New("invalid start date format")
	}

	endDate, err := ParseDate(endDateStr)
	if err != nil {
		return errors.New("invalid end date format")
	}

	if endDate.Before(startDate) {
		return errors.New("end date must be after start date")
	}

	return nil
}

func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	// Convert the first character to uppercase and append the rest of the string
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ToSnakeCase(s string) string {
	// Regular expression to find boundaries between lower and upper case letters
	re := regexp.MustCompile("([a-z])([A-Z])")
	// Replace those boundaries with an underscore followed by the uppercase letter in lowercase
	snake := re.ReplaceAllString(s, "${1}_${2}")
	// Convert the entire string to lowercase
	return strings.ToLower(snake)
}

func HasNumeric(input string) bool {
	// Regular expression to check for any numeric character
	re := regexp.MustCompile(`[0-9]`)
	return re.MatchString(input)
}

func IntsToString(nums []int) string {
	// Convert each int to a string
	strNums := make([]string, len(nums))
	for i, num := range nums {
		strNums[i] = strconv.Itoa(num)
	}
	// Join the string representations with commas (or any separator you prefer)
	return strings.Join(strNums, ",")
}

func GetUserEmails(users []models.User) []string {
	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}
	return emails
}

func GetUserUsernames(users []models.User) []string {
	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}
	return usernames
}
