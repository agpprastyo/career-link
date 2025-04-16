package mail

import (
	"context"
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// Client wraps SendGrid client
type Client struct {
	client *sendgrid.Client
	config *config.SendGridConfig
	log    *logger.Logger
}

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To           string
	Subject      string
	Body         string
	IsHTML       bool
	TemplateName string
	TemplateData map[string]interface{}
}

// NewSendGridClient creates a new SendGrid client
func NewSendGridClient(cfg *config.AppConfig, log *logger.Logger) *Client {
	if cfg.SendGrid.APIKey == "" {
		log.Warn("SENDGRID_API_KEY is not set")
	}

	client := sendgrid.NewSendClient(cfg.SendGrid.APIKey)

	return &Client{
		client: client,
		config: &cfg.SendGrid,
		log:    log,
	}
}

// SendEmail sends an email using SendGrid
func (c *Client) SendEmail(ctx context.Context, message EmailMessage) error {
	from := mail.NewEmail(c.config.FromName, c.config.FromEmail)
	to := mail.NewEmail("", message.To)

	content := mail.NewContent("text/plain", message.Body)
	if message.IsHTML {
		content = mail.NewContent("text/html", message.Body)
	}

	msg := mail.NewV3MailInit(from, message.Subject, to, content)

	response, err := c.client.Send(msg)
	if err != nil {
		c.log.WithError(err).Error("Failed to send email")
		return err
	}

	c.log.WithFields(logrus.Fields{
		"from":     c.config.FromEmail,
		"to":       message.To,
		"subject":  message.Subject,
		"status":   response.StatusCode,
		"response": response.Body,
	}).Info("Email sent successfully")

	return nil
}

// Close is a no-op for SendGrid but kept for interface consistency
func (c *Client) Close() {
	c.log.Info("Closing SendGrid client")
}
