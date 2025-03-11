package validator

import (
	"fmt"

	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

var (
	// RgxPhoneBasic Basic international phone number pattern
	RgxPhoneBasic = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

	// RgxPhoneIndonesia Indonesian phone number pattern (more specific)
	RgxPhoneIndonesia = regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)

	// RgxPhoneInternational International phone number pattern (more comprehensive)
	RgxPhoneInternational = regexp.MustCompile(`^\+(?:[0-9]●?){6,14}[0-9]$`)
)

// IsPhoneNumber validates a phone number using basic format
func IsPhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}
	return RgxPhoneBasic.MatchString(phone)
}

// IsIndonesianPhoneNumber validates specifically Indonesian phone numbers
func IsIndonesianPhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}
	return RgxPhoneIndonesia.MatchString(phone)
}

// IsInternationalPhoneNumber validates international phone numbers
func IsInternationalPhoneNumber(phone string) bool {
	if phone == "" {
		return false
	}
	return RgxPhoneInternational.MatchString(phone)
}

// FormatIndonesianPhoneNumber formats a phone number to Indonesian standard
func FormatIndonesianPhoneNumber(phone string) string {
	// Remove any non-digit characters
	digitsOnly := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// If the number starts with 0, replace it with "+62"
	if strings.HasPrefix(digitsOnly, "0") {
		digitsOnly = "62" + digitsOnly[1:]
	}

	// If the number doesn't start with 62, add it
	if !strings.HasPrefix(digitsOnly, "62") {
		digitsOnly = "62" + digitsOnly
	}

	return "+" + digitsOnly
}

// CleanPhoneNumber removes all non-digit characters except '+'
func CleanPhoneNumber(phone string) string {
	// Keep only digits and '+'
	return regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
}

// ValidateAndFormatPhone validates and formats a phone number
func ValidateAndFormatPhone(phone string) (string, error) {
	// Clean the phone number first
	cleaned := CleanPhoneNumber(phone)

	// Check if it's empty after cleaning
	if cleaned == "" {
		return "", fmt.Errorf("phone number is empty")
	}

	// Validate the cleaned number
	if !IsPhoneNumber(cleaned) {
		return "", fmt.Errorf("invalid phone number format")
	}

	// If it's an Indonesian number, format it accordingly
	if IsIndonesianPhoneNumber(cleaned) {
		return FormatIndonesianPhoneNumber(cleaned), nil
	}

	// Return the cleaned number for other cases
	return cleaned, nil
}

var (
	RgxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func Between[T constraints.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func In[T comparable](value T, safelist ...T) bool {
	for i := range safelist {
		if value == safelist[i] {
			return true
		}
	}
	return false
}

func AllIn[T comparable](values []T, safelist ...T) bool {
	for i := range values {
		if !In(values[i], safelist...) {
			return false
		}
	}
	return true
}

func NotIn[T comparable](value T, blocklist ...T) bool {
	for i := range blocklist {
		if value == blocklist[i] {
			return false
		}
	}
	return true
}

func NoDuplicates[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func IsEmail(value string) bool {
	if len(value) > 254 {
		return false
	}

	return RgxEmail.MatchString(value)
}

func IsURL(value string) bool {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

func IsNotAdminInput(value string) bool {
	return value != string(database.Adm)
}

func IsNotSuperAdminInput(value string) bool {
	return value != string(database.AdminRoleSuper)
}
