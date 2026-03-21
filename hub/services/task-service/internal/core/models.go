package core

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
	Port     int    `json:"port"`
}
