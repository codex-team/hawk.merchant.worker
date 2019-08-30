package main

type NotificationMessage struct {
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workspaceId"`
	Amount      uint32 `json:"amount"`
	Timestamp   uint32 `json:"timestamp"`
}
