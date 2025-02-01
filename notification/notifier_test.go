package notification

import (
	"testing"
)

// MockNotifier is a mock implementation of the Notifier interface for testing.
type MockNotifier struct {
	SentMessages []string
}

func (m *MockNotifier) Notify(userID, message string) error {
	m.SentMessages = append(m.SentMessages, message)
	return nil
}

func TestEmailNotifier_Notify(t *testing.T) {
	notifier := &EmailNotifier{}
	userID := "user123"
	message := "test message"

	// Simply call Notify and ensure it doesn't panic
	notifier.Notify(userID, message)

	// Since EmailNotifier just prints to the console, we can't easily verify the output.
	// This test is mostly to ensure the method doesn't panic.
}

func TestMockNotifier_Notify(t *testing.T) {
	mockNotifier := &MockNotifier{}
	userID := "user123"
	message := "test message"

	mockNotifier.Notify(userID, message)

	if len(mockNotifier.SentMessages) != 1 {
		t.Errorf("expected 1 message to be sent, got %d", len(mockNotifier.SentMessages))
	}

	if mockNotifier.SentMessages[0] != message {
		t.Errorf("expected message %q, got %q", message, mockNotifier.SentMessages[0])
	}
}
