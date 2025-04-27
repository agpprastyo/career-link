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
	// Ensure baseURL is properly formatted
	baseURL := r.verifyBaseURL + "/verify"
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	// Remove the path segment that's causing duplication
	// Change this line to use the correct path structure
	verificationURL := fmt.Sprintf("%s?email=%s&token=%s",
		baseURL,
		url.QueryEscape(email),
		url.QueryEscape(token))

	// Rest of the function remains unchanged
	data := templates.RegistrationData{
		Username:        username,
		Email:           email,
		Token:           token,
		VerificationURL: verificationURL,
		CompanyName:     "Career Link",
		SupportEmail:    "support@careerlink.com",
	}

	htmlContent, err := templates.GetRegistrationEmailHTML(data)
	if err != nil {
		return err
	}

	message := mail.EmailMessage{
		To:      email,
		Subject: "Verify Your Career Link Account",
		Body:    htmlContent,
		IsHTML:  true,
	}

	return r.mail.SendEmail(ctx, message)
}

// SendPostVerificationEmail Complete the SendPostVerificationEmail method in internal/user/repository/mail.go
func (r *UserRepository) SendPostVerificationEmail(ctx context.Context, username, email string) error {
	// Prepare template data
	loginURL := r.verifyBaseURL
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
