package main

import (
	"fmt"
	"time"

	"github.com/rafaelpolonio92/notification-service/notification"
)

func main() {
	rules := map[string]notification.RateLimit{
		"status":    {Count: 2, Window: time.Minute},
		"news":      {Count: 1, Window: 24 * time.Hour},
		"marketing": {Count: 3, Window: time.Hour},
	}

	emailNotifier := &notification.EmailNotifier{}
	service := notification.NewService(emailNotifier, rules)

	testCases := []struct {
		msgType, userID, message string
	}{
		{"news", "user", "news 1"},
		{"news", "user", "news 2"},
		{"news", "user", "news 3"},
		{"news", "another_user", "news 1"},
		{"update", "user", "update 1"},
	}

	for _, tc := range testCases {
		err := service.Send(tc.msgType, tc.userID, tc.message)
		if err != nil {
			fmt.Printf("Error sending %s to %s: %v\n", tc.msgType, tc.userID, err)
		}
	}
}
