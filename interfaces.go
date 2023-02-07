package connector

import "fmt"

// TokensPersistor mantains the tokens persistent
type TokensPersistor interface {
	GetAccessToken() string
	SetAccessToken(string)
}

// ResponseError with error codes
type ResponseError struct {
	ID     string `json:"error"`
	Reason string `json:"error_description"`
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("[%v] %v", err.ID, err.Reason)
}

// ResponseError.ID
const (
	ErrorInternal        = "internal"         // Connector error
	ErrorNetwork         = "network"          // Network error
	ErrorInvalidResponse = "invalid_response" // Server response cannot be decoded
	ErrorUnauthorized    = "unauthorized"     // Auth error
)

type authRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`

	ErrorID     string `json:"error"`
	ErrorReason string `json:"error_description"`
}
