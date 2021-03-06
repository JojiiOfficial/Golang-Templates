package models

//CredentialRequest request containing credentials
type CredentialRequest struct {
	Username string `json:"username"`
	Password string `json:"pass"`
}

//PingRequest ping request
type PingRequest struct {
	Payload string
}
