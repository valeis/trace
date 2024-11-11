package chat

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"io"
	"net/http"
)

type ChatRepository struct {
	cfg *config.Config
}

type ChatServiceResponse struct {
	Status bool          `json:"status"`
	Data   []interface{} `json:"data"`
	Total  int           `json:"total"`
}

type ChatHistory struct {
	Status  bool          `json:"status"`
	Total   int           `json:"total"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type ChatParticipants struct {
	User1 string `json:"u1"`
	User2 string `json:"u2"`
}

func NewChatRepository(cfg *config.Config) *ChatRepository {
	return &ChatRepository{cfg: cfg}
}

func (repo *ChatRepository) GetContactList(username string) (ChatServiceResponse, error) {
	userData := model.Username{username}
	usernamePayload, err := json.Marshal(userData)

	if err != nil {
		return ChatServiceResponse{}, fmt.Errorf("failed to marshal username: %v", err)
	}

	resp, err := http.Post(repo.cfg.ChatServiceURL+"contact-list", "application/json", bytes.NewBuffer(usernamePayload))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatServiceResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var chatResponse ChatServiceResponse
	if err := json.Unmarshal(responseBody, &chatResponse); err != nil {
		return ChatServiceResponse{}, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	if len(chatResponse.Data) <= 0 {
		slog.Errorf("Chat service returned an error: %s", chatResponse.Status)
		return ChatServiceResponse{}, fmt.Errorf("chat service error: %s", chatResponse.Status)
	}
	return chatResponse, nil
}

func (repo *ChatRepository) GetChatHistory(partcipants ChatParticipants) (ChatHistory, error) {
	participantsPayload, err := json.Marshal(partcipants)

	if err != nil {
		return ChatHistory{}, fmt.Errorf("failed to marshal participants: %v", err)
	}

	resp, err := http.Post(repo.cfg.ChatServiceURL+"chat-history", "application/json", bytes.NewBuffer(participantsPayload))

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatHistory{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var chatHistoryResponse ChatHistory
	if err := json.Unmarshal(responseBody, &chatHistoryResponse); err != nil {
		return ChatHistory{}, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	if len(chatHistoryResponse.Data) <= 0 {
		slog.Errorf("Chat history cannot be fetched: %s", chatHistoryResponse.Status)
	}
	return chatHistoryResponse, nil
}
