package templates

import (
	"bytes"
	_ "embed"
	"html/template"
	"time"
)

//go:embed html/post_verification.html
var postVerificationTemplate string

type PostVerificationData struct {
	Username     string
	Email        string
	LoginURL     string
	Year         int
	SupportEmail string
}

// GetPostVerificationHTML renders the post-verification email template
func GetPostVerificationHTML(data PostVerificationData) (string, error) {
	tmpl, err := template.New("post_verification").Parse(postVerificationTemplate)
	if err != nil {
		return "", err
	}

	// Set current year if not provided
	if data.Year == 0 {
		data.Year = time.Now().Year()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
