package templates

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed html/registration.html
var registrationEmailTemplate string

// RegistrationData contains the data needed for registration email
type RegistrationData struct {
	Username        string
	Email           string
	Token           string // Unique verification token for the link
	VerificationURL string // Base URL for the verification endpoint
	CompanyName     string
	SupportEmail    string
}

// GetRegistrationEmailHTML returns the HTML for registration verification email
func GetRegistrationEmailHTML(data RegistrationData) (string, error) {
	tmpl, err := template.New("registration").Parse(registrationEmailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
