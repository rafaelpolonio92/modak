package notification

import "fmt"

type Notifier interface {
	Notify(userID string, message string) error
}

type EmailNotifier struct{}

func (e *EmailNotifier) Notify(userID string, message string) error {
	fmt.Printf("EmailNotifier: sending email to user %s: %s\n", userID, message)
	return nil
}
