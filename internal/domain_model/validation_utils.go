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
	if len(*password) < 8 || len(*password) > 60 {
		return "Password length must be between 8 and 60"
	}
	return ""
}

func ValidateId(id int32) string {
	if id <= 0 {
		return "Invalid id"
	}
	return ""
}

func ValidateIds(ids []int32, assertAtLeastOne bool) map[string]string {
	messages := map[string]string{}
	if assertAtLeastOne && len(ids) == 0 {
		messages["ids"] = "At least one id is required"
		return messages
	}
	for _, id := range ids {
		msg := ValidateId(id)
		if msg != "" {
			messages[string(id)] = msg
		}
	}
	return messages
}

func ValidateRoleName(roleName string) string {
	if roleName == "" {
		return "Role name cannot be empty"
	}
	if len(roleName) < 1 || len(roleName) > 255 {
		return "Role name length must be between 1 and 255"
	}
	regexPattern := "^[A-Z0-9_]+$"
	matched, _ := regexp.MatchString(regexPattern, roleName)
	if !matched {
		return "Role name can only contain uppercase letters, numbers, and underscores"
	}
	return ""
}

func ValidateRoleNames(roleNames []string, assertAtLeastOne bool) map[string]string {
	messages := map[string]string{}
	if assertAtLeastOne && len(roleNames) == 0 {
		messages["role_names"] = "At least one role name is required"
		return messages
	}
	for _, roleName := range roleNames {
		msg := ValidateRoleName(roleName)
		if msg != "" {
			messages[roleName] = msg
		}
	}
	return messages
}
