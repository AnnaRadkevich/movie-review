package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/cloudmachinery/movie-reviews/internal/modules/users"

	"gopkg.in/validator.v2"
)

func SetupValidators() {
	validators := []struct {
		name string
		fn   validator.ValidationFunc
	}{
		{"password", password},
		{"email", email},
		{"role", role},
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

var (
	passwordMinLength       = 8
	emailMaxLength          = 127
	passwordSpecialChars    = "!$#()[]{}?+*~@^&-_"
	passwordRequiredEntries = []struct {
		name  string
		chars string
	}{
		{"lowercase character", "abcdefghijklmnopqrstuvwxyz"},
		{"uppercase character", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"digit", "0123456789"},
		{"special character (" + passwordSpecialChars + ")", passwordSpecialChars},
	}
)

func password(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("password only validates strings")
	}
	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
	}
	for _, required := range passwordRequiredEntries {
		if !strings.ContainsAny(s, required.chars) {
			return fmt.Errorf("password must contain at leats one %s", required.name)
		}
	}
	return nil
}

func email(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email only validates strings")
	}
	if len(s) > emailMaxLength {
		return fmt.Errorf("email must be at most %d characters long", emailMaxLength)
	}
	_, err := mail.ParseAddress(s)
	return err
}

func role(v interface{}, _ string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("role only validates string ")
	}
	if !(s == users.UserRole || s == users.EditorRole || s == users.AdminRole) {
		return fmt.Errorf("role must be only user/editor/admin: %d", s)
	}
	return nil
}
