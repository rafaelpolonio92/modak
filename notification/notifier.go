package notification

import "fmt"

type Notifier interface {
	Notify(userID int, message string) error
}

type EmailNotifier struct{}

func (e *EmailNotifier) Notify(userID int, message string) error {
	fmt.Printf("EmailNotifier: sending email to user %s: %s\n", userID, message)
	return nil
}
