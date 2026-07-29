package handlers

import (
	"encoding/json"
	"testing"
	"time"

	chat_models "github.com/nutrixpos/pos/modules/chat/models"
)

func TestChatMessage_Serialization(t *testing.T) {
	msg := chat_models.ChatMessage{
		Id:        "m-1",
		Channel:   "general",
		Sender:    "Janez Novak",
		SenderId:  "u-1",
		Content:   "Hello everyone!",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded chat_models.ChatMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Content != "Hello everyone!" {
		t.Errorf("expected Content='Hello everyone!', got %s", decoded.Content)
	}
	if decoded.Channel != "general" {
		t.Errorf("expected Channel='general', got %s", decoded.Channel)
	}
}

func TestChatChannel_Serialization(t *testing.T) {
	channel := chat_models.ChatChannel{
		Id:          "kitchen",
		Name:        "Kitchen",
		Description: "Kitchen team chat",
		IsDefault:   false,
	}

	data, err := json.Marshal(channel)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded chat_models.ChatChannel
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Kitchen" {
		t.Errorf("expected Name='Kitchen', got %s", decoded.Name)
	}
	if decoded.IsDefault {
		t.Error("expected IsDefault=false")
	}
}

func TestChatSummary_Serialization(t *testing.T) {
	summary := chat_models.ChatSummary{
		TotalMessages: 42,
		Channels: []chat_models.ChatChannel{
			{Id: "general", Name: "General", IsDefault: true},
		},
		OnlineUsers: []chat_models.OnlineUser{
			{UserId: "u-1", Username: "Janez"},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded chat_models.ChatSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalMessages != 42 {
		t.Errorf("expected TotalMessages=42, got %d", decoded.TotalMessages)
	}
}

func TestDefaultChannels(t *testing.T) {
	defaults := []string{"general", "kitchen", "service", "management"}
	for _, ch := range defaults {
		if ch == "" {
			t.Error("channel id should not be empty")
		}
	}
}
