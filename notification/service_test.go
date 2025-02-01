package notification

import (
	"testing"
	"time"
)

func TestService_Send(t *testing.T) {
	rules := map[string]RateLimit{
		"status": {Count: 2, Window: time.Minute},
		"news":   {Count: 1, Window: 24 * time.Hour},
	}

	mockNotifier := &MockNotifier{}
	service := NewService(mockNotifier, rules)

	err := service.Send("status", "user1", "status update 1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.Send("status", "user1", "status update 2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.Send("status", "user1", "status update 3")
	if err != ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}

	err = service.Send("status", "user2", "status update 1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.Send("news", "user1", "news update 1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.Send("news", "user1", "news update 2")
	if err != ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
}

func TestService_Send_TimeWindow(t *testing.T) {
	rules := map[string]RateLimit{
		"status": {Count: 2, Window: time.Minute},
	}

	mockNotifier := &MockNotifier{}
	service := NewService(mockNotifier, rules)

	// Set a fixed clock for testing
	fixedTime := time.Now()
	mockClock := &MockClock{now: fixedTime}
	service.SetClock(mockClock)

	// Send two notifications within the time window
	err := service.Send("status", "user1", "status update 1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.Send("status", "user1", "status update 2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Exceed the rate limit
	err = service.Send("status", "user1", "status update 3")
	if err != ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}

	// Advance time beyond the window and try again
	mockClock.now = fixedTime.Add(time.Minute + time.Second)
	err = service.Send("status", "user1", "status update 3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// MockClock is a mock implementation of the Clock interface for testing.
type MockClock struct {
	now time.Time
}

func (m *MockClock) Now() time.Time {
	return m.now
}
