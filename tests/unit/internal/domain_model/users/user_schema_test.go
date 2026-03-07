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
				Roles:     []string{"ADMIN", "TEACHER"},
			},
			expectError: false,
		},
		{
			name: "valid create user with single role",
			input: &users.CreateUser{
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@example.com",
				Password:  "password123",
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"ADMIN"},
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
				Roles:     []string{"TEACHER"},
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
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"ADMIN"},
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
				Roles:     []string{"TEACHER"},
			},
			expectError: true,
			expectData: map[string]string{
				"password": "Password length must be between 8 and 60",
			},
		},
		{
			name: "missing roles",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
				Roles:     []string{},
			},
			expectError: true,
			expectData: map[string]string{
				"roles": "At least one role is required to create a user",
			},
		},
		{
			name: "invalid role name - lowercase",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
				Roles:     []string{"admin"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[admin]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "invalid role name - special characters",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
				Roles:     []string{"ADMIN-ROLE"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[ADMIN-ROLE]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "multiple invalid role names",
			input: &users.CreateUser{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Password:  "password123",
				Roles:     []string{"admin", "Teacher", "STUDENT-1"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[admin]":     "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[Teacher]":   "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[STUDENT-1]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "multiple validation errors",
			input: &users.CreateUser{
				FirstName: "",
				LastName:  "",
				Email:     "bad",
				Password:  "",
				Roles:     []string{},
			},
			expectError: true,
			expectData: map[string]string{
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
				"password":   "Password cannot be empty",
				"roles":      "At least one role is required to create a user",
			},
		},
		{
			name: "multiple validation errors including invalid roles",
			input: &users.CreateUser{
				FirstName: "",
				LastName:  "",
				Email:     "bad",
				Password:  "short",
				Roles:     []string{"invalid-role", "lowercase"},
			},
			expectError: true,
			expectData: map[string]string{
				"first_name":          "Field cannot be empty",
				"last_name":           "Field cannot be empty",
				"email":               "Invalid email format",
				"password":            "Password length must be between 8 and 60",
				"roles[invalid-role]": "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[lowercase]":    "Role name can only contain uppercase letters, numbers, and underscores",
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

			assert.Equal(t, "VALIDATION_ERROR", err.Code)
			assert.Equal(t, "Validation failed during user creation", err.Message)
			assert.Equal(t, tt.expectData, err.Data)
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
				Roles:     []string{"ADMIN", "TEACHER"},
			},
			expectError: false,
		},
		{
			name: "valid update user with single role",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "John",
				LastName:  "Smith",
				Email:     "john.smith@example.com",
				Roles:     []string{"STUDENT"},
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
				Roles:     []string{"ADMIN"},
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
				Roles:     []string{"TEACHER"},
			},
			expectError: true,
			expectData: map[string]string{
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
			},
		},
		{
			name: "missing roles",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Roles:     []string{},
			},
			expectError: true,
			expectData: map[string]string{
				"roles": "At least one role is required to update a user",
			},
		},
		{
			name: "invalid role name - lowercase",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Roles:     []string{"student"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[student]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "invalid role name - mixed case",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Roles:     []string{"Admin", "Teacher"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[Admin]":   "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[Teacher]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "invalid role name - special characters",
			input: &users.UpdateUser{
				UserId:    1,
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane.doe@example.com",
				Roles:     []string{"ROLE-ADMIN", "ROLE@TEACHER"},
			},
			expectError: true,
			expectData: map[string]string{
				"roles[ROLE-ADMIN]":   "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[ROLE@TEACHER]": "Role name can only contain uppercase letters, numbers, and underscores",
			},
		},
		{
			name: "multiple validation errors",
			input: &users.UpdateUser{
				UserId:    -1,
				FirstName: "",
				LastName:  "",
				Email:     "bad",
				Roles:     []string{},
			},
			expectError: true,
			expectData: map[string]string{
				"user_id":    "Invalid id",
				"first_name": "Field cannot be empty",
				"last_name":  "Field cannot be empty",
				"email":      "Invalid email format",
				"roles":      "At least one role is required to update a user",
			},
		},
		{
			name: "multiple validation errors including invalid roles",
			input: &users.UpdateUser{
				UserId:    0,
				FirstName: "",
				LastName:  "",
				Email:     "bad",
				Roles:     []string{"admin", "Teacher-1"},
			},
			expectError: true,
			expectData: map[string]string{
				"user_id":          "Invalid id",
				"first_name":       "Field cannot be empty",
				"last_name":        "Field cannot be empty",
				"email":            "Invalid email format",
				"roles[admin]":     "Role name can only contain uppercase letters, numbers, and underscores",
				"roles[Teacher-1]": "Role name can only contain uppercase letters, numbers, and underscores",
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
				"password": "Password length must be between 8 and 60",
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
