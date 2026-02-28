package domain_model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/tests"
)

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
			input:    tests.StrPtr("hello"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "valid string at minimum length",
			input:    tests.StrPtr("abc"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "valid string at maximum length",
			input:    tests.StrPtr("abcdefghij"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "",
		},
		{
			name:     "empty string when required",
			input:    tests.StrPtr(""),
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
			input:    tests.StrPtr(""),
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
			input:    tests.StrPtr("ab"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field length must be between 3 and 10",
		},
		{
			name:     "string too long",
			input:    tests.StrPtr("abcdefghijk"),
			minLen:   3,
			maxLen:   10,
			required: true,
			want:     "Field length must be between 3 and 10",
		},
		{
			name:     "single character with min 1",
			input:    tests.StrPtr("a"),
			minLen:   1,
			maxLen:   5,
			required: true,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model.ValidateString(tt.input, tt.minLen, tt.maxLen, tt.required)
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
			email:    tests.StrPtr("test@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with dots",
			email:    tests.StrPtr("user.name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with hyphens",
			email:    tests.StrPtr("user-name@example-domain.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with underscore",
			email:    tests.StrPtr("user_name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - with numbers",
			email:    tests.StrPtr("user123@example123.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - subdomain",
			email:    tests.StrPtr("user@mail.example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - multiple subdomains",
			email:    tests.StrPtr("user@mail.sub.example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "valid email - two letter TLD",
			email:    tests.StrPtr("user@example.io"),
			required: true,
			want:     "",
		},
		{
			name:     "empty email when required",
			email:    tests.StrPtr(""),
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
			email:    tests.StrPtr(""),
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
			email:    tests.StrPtr("userexample.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no domain",
			email:    tests.StrPtr("user@"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no local part",
			email:    tests.StrPtr("@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - no TLD",
			email:    tests.StrPtr("user@example"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - multiple @ symbols",
			email:    tests.StrPtr("user@@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - spaces",
			email:    tests.StrPtr("user name@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "invalid email - special characters",
			email:    tests.StrPtr("user#name@example.com"),
			required: true,
			want:     "Invalid email format",
		},
		{
			name:     "email starts with dot - currently passes regex",
			email:    tests.StrPtr(".user@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "email ends with dot - currently passes regex",
			email:    tests.StrPtr("user.@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "email with consecutive dots - currently passes regex",
			email:    tests.StrPtr("user..name@example.com"),
			required: true,
			want:     "",
		},
		{
			name:     "invalid email - single letter TLD",
			email:    tests.StrPtr("user@example.c"),
			required: true,
			want:     "Invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model.ValidateEmail(tt.email, tt.required)
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
			password: tests.StrPtr("password"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - 32 characters",
			password: tests.StrPtr("12345678901234567890123456789012"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - middle range",
			password: tests.StrPtr("mySecurePassword123"),
			required: true,
			want:     "",
		},
		{
			name:     "valid password - with special characters",
			password: tests.StrPtr("P@ssw0rd!"),
			required: true,
			want:     "",
		},
		{
			name:     "empty password when required",
			password: tests.StrPtr(""),
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
			password: tests.StrPtr(""),
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
			password: tests.StrPtr("passwor"),
			required: true,
			want:     "Password length must be between 8 and 60",
		},
		{
			name:     "password too short - 1 character",
			password: tests.StrPtr("p"),
			required: true,
			want:     "Password length must be between 8 and 60",
		},
		{
			name:     "password too long - 61 characters",
			password: tests.StrPtr("1234567890123456789012345678901231234567890123456789012345678"),
			required: true,
			want:     "Password length must be between 8 and 60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain_model.ValidatePassword(tt.password, tt.required)
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
			got := domain_model.ValidateId(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}
