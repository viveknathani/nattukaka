package utils

import (
	"math/rand"
	"net/http"
	"os"
	"regexp"

	"github.com/mattevans/postmark-go"
)

// IsValidEmail checks if an email address is valid using a regex.
func IsValidEmail(email string) bool {
	// Regular expression for basic email validation.
	// This is a simple pattern and might not cover all edge cases.
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// GenerateRandomNumber generates a random integer between min (inclusive) and max (exclusive).
func GenerateRandomNumber(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

// SendEmail sends an email via Postmark
func SendEmail(to string, subject string, body string) error {
	// Init client with round tripper adding auth fields.
	client := postmark.NewClient(
		postmark.WithClient(&http.Client{
			Transport: &postmark.AuthTransport{Token: os.Getenv("POSTMARK_API_KEY")},
		}),
	)

	// Build the email.
	emailReq := &postmark.Email{
		From:     "vivek@teachyourselfmath.app",
		To:       to,
		Subject:  subject,
		HTMLBody: body,
	}

	// Send it!
	_, _, err := client.Email.Send(emailReq)
	return err
}
