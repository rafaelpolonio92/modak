package notification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockClock struct {
	now time.Time
}

func (m *MockClock) Now() time.Time {
	return m.now
}

func TestService_Send(t *testing.T) {
	rules := map[string]RateLimit{
		"status": {Count: 2, Window: time.Minute},
		"news":   {Count: 1, Window: 24 * time.Hour},
	}

	mockNotifier := &MockNotifier{}
	service := NewService(mockNotifier, rules)

	t.Run("Status Notifications within limit", func(t *testing.T) {
		assert.NoError(t, service.Send("status", "user1", "status update 1"))
		assert.NoError(t, service.Send("status", "user1", "status update 2"))
	})

	t.Run("Status Notifications exceeding limit", func(t *testing.T) {
		err := service.Send("status", "user1", "status update 3")
		assert.ErrorIs(t, err, ErrRateLimitExceeded)
	})

	t.Run("Different User Status Notification", func(t *testing.T) {
		assert.NoError(t, service.Send("status", "user2", "status update 1"))
	})

	t.Run("News Notification within limit", func(t *testing.T) {
		assert.NoError(t, service.Send("news", "user1", "news update 1"))
	})

	t.Run("News Notification exceeding limit", func(t *testing.T) {
		err := service.Send("news", "user1", "news update 2")
		assert.ErrorIs(t, err, ErrRateLimitExceeded)
	})
}

func TestService_Send_TimeWindow(t *testing.T) {
	rules := map[string]RateLimit{
		"status": {Count: 2, Window: time.Minute},
	}

	mockNotifier := &MockNotifier{}
	service := NewService(mockNotifier, rules)

	fixedTime := time.Now()
	mockClock := &MockClock{now: fixedTime}
	service.SetClock(mockClock)

	t.Run("Status Notifications within limit", func(t *testing.T) {
		assert.NoError(t, service.Send("status", "user1", "status update 1"))
		assert.NoError(t, service.Send("status", "user1", "status update 2"))
	})

	t.Run("Status Notification exceeding limit", func(t *testing.T) {
		err := service.Send("status", "user1", "status update 3")
		assert.ErrorIs(t, err, ErrRateLimitExceeded)
	})

	t.Run("Reset after time window", func(t *testing.T) {
		mockClock.now = fixedTime.Add(time.Minute + time.Second)
		assert.NoError(t, service.Send("status", "user1", "status update 3"))
	})
}
