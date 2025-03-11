package validator

import (
	"regexp"
	"strings"
)

type Validator struct {
	Errors      []string          `json:",omitempty"`
	FieldErrors map[string]string `json:",omitempty"`
}

func (v Validator) HasErrors() bool {
	return len(v.Errors) != 0 || len(v.FieldErrors) != 0
}

func (v *Validator) AddError(message string) {
	if v.Errors == nil {
		v.Errors = []string{}
	}

	v.Errors = append(v.Errors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func ValidatePhone(phone string) (string, []string) {
	var errors []string

	// If phone is empty, return early
	if phone == "" {
		return phone, nil
	}

	// Check length
	if !MaxRunes(phone, 20) {
		errors = append(errors, "must be less than 20 characters")
	}

	// Clean the phone number
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Check basic format
	if !IsPhoneNumber(cleanPhone) {
		errors = append(errors, "must be a valid phone number format")
	}

	// Check Indonesian format if applicable
	if strings.HasPrefix(cleanPhone, "+62") || strings.HasPrefix(cleanPhone, "62") || strings.HasPrefix(cleanPhone, "0") {
		if !IsIndonesianPhoneNumber(cleanPhone) {
			errors = append(errors, "must be a valid Indonesian phone number")
		}
	}

	return cleanPhone, errors
}
