package responses

import "github.com/ong-gtp/go-chat/models"

type SignUpResponse struct {
	User  models.User `json:"User"`
	Token string      `json:"Token"`
}
