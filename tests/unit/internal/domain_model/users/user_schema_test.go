package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todorpopov/school-manager/internal/domain_model/users"
)

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		input       *users.CreateUser
		expectError bool
		expectData  map[string]string
	}{
		{
			name: "valid create user",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
			},
			expectError: false,
		},
		{
			name: "missing first name",
			input: &users.CreateUser{
				FirstName: "",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
			},
			expectError: true,
			expectData: map[string]string{
				"first_name": "Field cannot be empty",
			},
		},
		{
			name: "missing last name",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "",
				Email:     "jane.doe@example.com",
				Password:  "password123",
			},
			expectError: true,
			expectData: map[string]string{
				"last_name": "Field cannot be empty",
			},
		},
		{
			name: "invalid email",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "password123",
			},
			expectError: true,
			expectData: map[string]string{
				"email": "Invalid email format",
			},
		},
		{
			name: "missing password",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "",
			},
			expectError: true,
			expectData: map[string]string{
				"password": "Password cannot be empty",
			},
		},
		{
			name: "short password",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "short",
			},
			expectError: true,
			expectData: map[string]string{
				"password": "Password length must be between 8 and 32",
			},
		},
		{
			name: "multiple validation errors",
			input: &users.CreateUser{
				FirstName: "",
				LastName:  "",
				Email:     "bad",
				Password:  "",
			},
			expectError: true,
			expectData: map[string]string{
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
				"password":   "Password cannot be empty",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := users.ValidateCreateUser(tt.input)
			if !tt.expectError {
				assert.Nil(t, err)
				return
			}

			if assert.NotNil(t, err) {
				assert.Equal(t, "VALIDATION_ERROR", err.Code)
				assert.Equal(t, "Validation failed during user creation", err.Message)
				assert.Equal(t, tt.expectData, err.Data)
			}
		})
	}
}

func TestValidateUpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		input       *users.UpdateUser
		expectError bool
		expectData  map[string]string
	}{
		{
			name: "valid update user",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
			},
			expectError: false,
		},
		{
			name: "invalid user id",
			input: &users.UpdateUser{
				UserId:    0,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
			},
			expectError: true,
			expectData: map[string]string{
				"user_id": "Invalid id",
			},
		},
		{
			name: "invalid names and email",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "",
				LastName:  "",
				Email:     "invalid-email",
			},
			expectError: true,
			expectData: map[string]string{
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
			},
		},
		{
			name: "multiple validation errors",
			input: &users.UpdateUser{
				UserId:    -1,
				FirstName: "",
				LastName:  "",
				Email:     "bad",
			},
			expectError: true,
			expectData: map[string]string{
				"user_id":    "Invalid id",
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := users.ValidateUpdateUser(tt.input)
			if !tt.expectError {
				assert.Nil(t, err)
				return
			}

			if assert.NotNil(t, err) {
				assert.Equal(t, "VALIDATION_ERROR", err.Code)
				assert.Equal(t, "Validation failed during user update", err.Message)
				assert.Equal(t, tt.expectData, err.Data)
			}
		})
	}
}

func TestValidateUpdateUserPassword(t *testing.T) {
	tests := []struct {
		name        string
		input       *users.UpdateUserPassword
		expectError bool
		expectData  map[string]string
	}{
		{
			name: "valid update user password",
			input: &users.UpdateUserPassword{
				UserId:   1,
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "invalid user id",
			input: &users.UpdateUserPassword{
				UserId:   0,
				Password: "password123",
			},
			expectError: true,
			expectData: map[string]string{
				"user_id": "Invalid id",
			},
		},
		{
			name: "invalid password",
			input: &users.UpdateUserPassword{
				UserId:   1,
				Password: "short",
			},
			expectError: true,
			expectData: map[string]string{
				"password": "Password length must be between 8 and 32",
			},
		},
		{
			name: "multiple validation errors",
			input: &users.UpdateUserPassword{
				UserId:   -2,
				Password: "",
			},
			expectError: true,
			expectData: map[string]string{
				"user_id":  "Invalid id",
				"password": "Password cannot be empty",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := users.ValidateUpdateUserPassword(tt.input)
			if !tt.expectError {
				assert.Nil(t, err)
				return
			}

			if assert.NotNil(t, err) {
				assert.Equal(t, "VALIDATION_ERROR", err.Code)
				assert.Equal(t, "Validation failed during user password update", err.Message)
				assert.Equal(t, tt.expectData, err.Data)
			}
		})
	}
}
