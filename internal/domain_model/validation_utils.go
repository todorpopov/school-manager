package domain_model

import (
	"fmt"
	"regexp"
)

func ValidateString(s *string, minLen int, maxLen int, required bool) string {
	if s == nil || *s == "" {
		if required == true {
			return "Field cannot be empty"
		}
		return ""
	}
	if len(*s) < minLen || len(*s) > maxLen {
		return fmt.Sprintf("Field length must be between %d and %d", minLen, maxLen)
	}
	return ""
}

func ValidateEmail(email *string, required bool) string {
	if email == nil || *email == "" {
		if required == true {
			return "Email cannot be empty"
		}
		return ""
	}

	regexPattern := "^[A-Za-z0-9._-]{1,64}@[A-Za-z0-9-]{1,63}(\\.[A-Za-z0-9-]{1,63})*\\.[A-Za-z]{2,}$"
	matched, _ := regexp.MatchString(regexPattern, *email)
	if !matched {
		return "Invalid email format"
	}
	return ""
}

func ValidatePassword(password *string, required bool) string {
	if password == nil || *password == "" {
		if required == true {
			return "Password cannot be empty"
		}
		return ""
	}
	if len(*password) < 8 || len(*password) > 32 {
		return "Password length must be between 8 and 32"
	}
	return ""
}

func ValidateId(id int32) string {
	if id <= 0 {
		return "Invalid id"
	}
	return ""
}
