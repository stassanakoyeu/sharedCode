package models

type EchoMessage struct {
	Message string `json:"message" binding:"required"`
}
