package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockNotifier struct {
	SentMessages []string
}

func (m *MockNotifier) Notify(userID, message string) error {
	m.SentMessages = append(m.SentMessages, message)
	return nil
}

func TestNotifier(t *testing.T) {
	t.Run("Test EmailNotifier", func(t *testing.T) {
		notifier := &EmailNotifier{}
		userID := "user123"
		message := "test message"

		err := notifier.Notify(userID, message)
		assert.NoError(t, err, "Notify should not return an error")
	})

	t.Run("Test MockNotifier", func(t *testing.T) {
		mockNotifier := &MockNotifier{}
		userID := "user123"
		message := "test message"

		err := mockNotifier.Notify(userID, message)
		assert.NoError(t, err, "Notify should not return an error")

		assert.Equal(t, 1, len(mockNotifier.SentMessages), "expected 1 message to be sent")
		assert.Equal(t, message, mockNotifier.SentMessages[0], "expected message to match")
	})
}
