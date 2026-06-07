package dto

type UserResponse struct {
	ID        string `json:"id"`
	DiscordID string `json:"discord_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
}

type MessageResponse struct {
	Message string `json:"message"` // generic message response for success or error messages
}