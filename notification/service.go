package notification

import (
	"errors"
	"sync"
	"time"
)

var ErrRateLimitExceeded = errors.New("rate limit exceeded")

type RateLimit struct {
	Count  int
	Window time.Duration
}

type Timer interface {
	Now() time.Time
}

type key struct {
	userID  string
	msgType string
}

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now()
}

type Service struct {
	notifier Notifier
	rules    map[string]RateLimit
	mu       sync.Mutex
	history  map[key][]time.Time
	timer    Timer
}

func NewService(notifier Notifier, rules map[string]RateLimit) *Service {
	return &Service{
		notifier: notifier,
		rules:    rules,
		history:  make(map[key][]time.Time),
		timer:    &Clock{},
	}
}

func (s *Service) SetClock(timer Timer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timer = timer
}

func (s *Service) Send(msgType, userID, message string) error {
	rule, exists := s.rules[msgType]
	if !exists {
		s.notifier.Notify(userID, message)
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	k := key{userID: userID, msgType: msgType}
	now := s.timer.Now()
	cutoff := now.Add(-rule.Window)

	timestamps := s.history[k]
	var valid []time.Time
	count := 0
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
			count++
		}
	}

	if count >= rule.Count {
		return ErrRateLimitExceeded
	}

	valid = append(valid, now)
	s.history[k] = valid

	s.notifier.Notify(userID, message)
	return nil
}
