package templates

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed html/verification_success.html
var successVerificationTemplate string

type GetVerificationSuccessData struct {
	LoginURL string
}

func GetVerificationSuccessHTML(loginURL string) (string, error) {
	tmpl, err := template.New("verification_success").Parse(successVerificationTemplate)
	if err != nil {
		return "", err
	}

	data := GetVerificationSuccessData{
		LoginURL: loginURL,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil

}
