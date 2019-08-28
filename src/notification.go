package main

type NotificationMessage struct {
	Message     string `json:"message"`
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workspaceId"`
	Amount      uint64 `json:"amount"`
	Timestamp   uint32 `json:"timestamp"`
}
