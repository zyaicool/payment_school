package utilities

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidateFieldNotEmpty(t *testing.T) {
	tests := []struct {
		field     interface{}
		fieldName string
		wantErr   error
	}{
		// Test for empty string
		{field: "", fieldName: "Name", wantErr: errors.New("Name cannot be empty")},
		// Test for valid string
		{field: "Valid", fieldName: "Name", wantErr: nil},
		// Test for zero int
		{field: 0, fieldName: "Age", wantErr: errors.New("Age cannot be zero or empty")},
		// Test for valid int
		{field: 25, fieldName: "Age", wantErr: nil},
		// Test for unsupported type
		{field: true, fieldName: "IsValid", wantErr: fmt.Errorf("IsValid has an unsupported type")},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing %s", tt.fieldName), func(t *testing.T) {
			err := ValidateFieldNotEmpty(tt.field, tt.fieldName)
			if err != nil && tt.wantErr == nil {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error %v, got nil", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
