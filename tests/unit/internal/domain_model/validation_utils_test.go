package domain_model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	domain_model2 "github.com/todorpopov/school-manager/internal/domain_model"
)

func strPtr(s string) *string {
	return &s
}

func TestValidateString(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		minLen   int
		maxLen   int
		required bool
		want     string
	}{
		{
			name:     "valid string within range",
			input:    strPtr("hello"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "valid string at minimum length",
			input:    strPtr("abc"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "valid string at maximum length",
			input:    strPtr("abcdefghij"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "empty string when required",
			input:    strPtr(""),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field cannot be empty",
		},
		{
			name:     "nil string when required",
			input:    nil,
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field cannot be empty",
		},
		{
			name:     "empty string when not required",
			input:    strPtr(""),
			minLen:   3,
			maxLen:   10,
			required: false,
			want:     "",
		},
		{
			name:     "nil string when not required",
			input:    nil,
			minLen:   3,
			maxLen:   10,
			required: false,
			want:     "",
		},
		{
			name:     "string too short",
			input:    strPtr("ab"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field length must be between 3 and 10",
		},
		{
			name:     "string too long",
			input:    strPtr("abcdefghijk"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field length must be between 3 and 10",
		},
		{
			name:     "single character with min 1",
			input:    strPtr("a"),
			minLen:   1,
			maxLen:   5,
			required: true,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model2.ValidateString(tt.input, tt.minLen, tt.maxLen, tt.required)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    *string
		required bool
		want     string
	}{
		{
			name:     "valid email - simple",
			email:    strPtr("test@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with dots",
			email:    strPtr("user.name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with hyphens",
			email:    strPtr("user-name@example-domain.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with underscore",
			email:    strPtr("user_name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with numbers",
			email:    strPtr("user123@example123.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - subdomain",
			email:    strPtr("user@mail.example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - multiple subdomains",
			email:    strPtr("user@mail.sub.example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - two letter TLD",
			email:    strPtr("user@example.io"),
			required: true,
			want:     "",
		},
		{
			name:     "empty email when required",
			email:    strPtr(""),
			required: true,
			want:     "Email cannot be empty",
		},
		{
			name:     "nil email when required",
			email:    nil,
			required: true,
			want:     "Email cannot be empty",
		},
		{
			name:     "empty email when not required",
			email:    strPtr(""),
			required: false,
			want:     "",
		},
		{
			name:     "nil email when not required",
			email:    nil,
			required: false,
			want:     "",
		},
		{
			name:     "invalid email - no @ symbol",
			email:    strPtr("userexample.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no domain",
			email:    strPtr("user@"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no local part",
			email:    strPtr("@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no TLD",
			email:    strPtr("user@example"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - multiple @ symbols",
			email:    strPtr("user@@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - spaces",
			email:    strPtr("user name@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - special characters",
			email:    strPtr("user#name@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "email starts with dot - currently passes regex",
			email:    strPtr(".user@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "email ends with dot - currently passes regex",
			email:    strPtr("user.@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "email with consecutive dots - currently passes regex",
			email:    strPtr("user..name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "invalid email - single letter TLD",
			email:    strPtr("user@example.c"),
			required: true,
			want:     "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model2.ValidateEmail(tt.email, tt.required)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password *string
		required bool
		want     string
	}{
		{
			name:     "valid password - 8 characters",
			password: strPtr("password"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - 32 characters",
			password: strPtr("12345678901234567890123456789012"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - middle range",
			password: strPtr("mySecurePassword123"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - with special characters",
			password: strPtr("P@ssw0rd!"),
			required: true,
			want:     "",
		},
		{
			name:     "empty password when required",
			password: strPtr(""),
			required: true,
			want:     "Password cannot be empty",
		},
		{
			name:     "nil password when required",
			password: nil,
			required: true,
			want:     "Password cannot be empty",
		},
		{
			name:     "empty password when not required",
			password: strPtr(""),
			required: false,
			want:     "",
		},
		{
			name:     "nil password when not required",
			password: nil,
			required: false,
			want:     "",
		},
		{
			name:     "password too short - 7 characters",
			password: strPtr("passwor"),
			required: true,
			want:     "Password length must be between 8 and 32",
		},
		{
			name:     "password too short - 1 character",
			password: strPtr("p"),
			required: true,
			want:     "Password length must be between 8 and 32",
		},
		{
			name:     "password too long - 33 characters",
			password: strPtr("123456789012345678901234567890123"),
			required: true,
			want:     "Password length must be between 8 and 32",
		},
		{
			name:     "password too long - 50 characters",
			password: strPtr("12345678901234567890123456789012345678901234567890"),
			required: true,
			want:     "Password length must be between 8 and 32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model2.ValidatePassword(tt.password, tt.required)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateId(t *testing.T) {
	tests := []struct {
		name string
		id   int32
		want string
	}{
		{
			name: "valid id - positive number",
			id:   1,
			want: "",
		},
		{
			name: "valid id - large positive number",
			id:   999999,
			want: "",
		},
		{
			name: "valid id - max int32",
			id:   2147483647,
			want: "",
		},
		{
			name: "invalid id - zero",
			id:   0,
			want: "Invalid id",
		},
		{
			name: "invalid id - negative number",
			id:   -1,
			want: "Invalid id",
		},
		{
			name: "invalid id - large negative number",
			id:   -999999,
			want: "Invalid id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model2.ValidateId(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}
