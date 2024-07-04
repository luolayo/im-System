package test

import (
	"im-System/model"
	"testing"
)

// TestUserEventGetType tests the GetType method of UserEvent class
func TestUserEventGetType(t *testing.T) {
	// Create a user object
	user := model.NewUser(&mockConn{}, "testuser")

	// Create a user event (join event)
	event := &model.UserEvent{
		Type: model.UserJoin,
		User: user,
	}

	// Get the event type
	eventType := event.Type

	// Check if the event type is correct
	if eventType != model.UserJoin {
		t.Errorf("Expected event type %s, got %s", model.UserJoin, eventType)
	}
}
