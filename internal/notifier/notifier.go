// internal/notifier/notifier.go
package notifier

import (
	"fmt"
)

type EmailNotifier struct {
	From    string
	To      string
	Subject string

	SMTPHost string
	SMTPPort int
}

type SmsNotifier struct {
	From        string
	To          string
	Provider    string
	BodyPreview string
}

var _ Notifier = (*EmailNotifier)(nil)
var _ Notifier = (*SmsNotifier)(nil)

func (e EmailNotifier) Send(message string) error {
	if !e.IsValid() {
		return fmt.Errorf("email invalide : %+v", e)
	}
	fmt.Printf("[Email envoyé à %s via %s:%d] %s\n", e.To, e.SMTPHost, e.SMTPPort, message)
	return nil
}

func (e EmailNotifier) IsValid() bool {
	return e.From != "" && e.To != "" && e.SMTPHost != "" && e.SMTPPort != 0
}

func (s SmsNotifier) Send(message string) error {
	if !s.IsValid() {
		return fmt.Errorf("SMS invalide : %+v", s)
	}
	fmt.Printf("[SMS envoyé à %s via %s] %s\n", s.To, s.Provider, message)
	return nil
}

func (s SmsNotifier) IsValid() bool {
	return s.From != "" && s.To != "" && s.Provider != ""
}
