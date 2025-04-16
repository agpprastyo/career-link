package repository

import (
	"context"
	"fmt"
	"github.com/agpprastyo/career-link/pkg/mail"
	"github.com/agpprastyo/career-link/pkg/mail/templates"
	"net/url"
	"strings"
	"time"
)

func (r *UserRepository) SendVerificationEmail(ctx context.Context, username, email, token string) error {

	baseURL := r.verifyBaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	verificationURL := fmt.Sprintf("%sverify?email=%s&token=%s",
		baseURL,
		url.QueryEscape(email),
		url.QueryEscape(token))

	// Prepare template data
	data := templates.RegistrationData{
		Username:        username,
		Email:           email,
		Token:           token,
		VerificationURL: verificationURL,
		CompanyName:     "Career Link",
		SupportEmail:    "support@careerlink.com",
	}

	// Generate HTML from template
	htmlContent, err := templates.GetRegistrationEmailHTML(data)
	if err != nil {
		return err
	}

	// Send email using SendGrid
	message := mail.EmailMessage{
		To:      email,
		Subject: "Verify Your Career Link Account",
		Body:    htmlContent,
		IsHTML:  true,
	}

	err = r.mail.SendEmail(ctx, message)
	return err
}

// SendPostVerificationEmail Complete the SendPostVerificationEmail method in internal/user/repository/mail.go
func (r *UserRepository) SendPostVerificationEmail(ctx context.Context, username, email string) error {
	// Prepare template data
	loginURL := strings.TrimSuffix(r.verifyBaseURL, "/verify")
	if !strings.HasSuffix(loginURL, "/") {
		loginURL += "/"
	}
	loginURL += "login"

	data := templates.PostVerificationData{
		Username:     username,
		Email:        email,
		LoginURL:     loginURL,
		Year:         time.Now().Year(),
		SupportEmail: "support@careerlink.com",
	}

	// Generate HTML from template
	htmlContent, err := templates.GetPostVerificationHTML(data)
	if err != nil {
		return err
	}

	// Send email
	message := mail.EmailMessage{
		To:      email,
		Subject: "Welcome to Career Link - Your Account is Active",
		Body:    htmlContent,
		IsHTML:  true,
	}

	return r.mail.SendEmail(ctx, message)
}
