package dto

import (
	"exaroton-wa-bot/internal/database/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserLoginReq struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func (r *UserLoginReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 20)),
	)
}

// Login page
type LoginPageData struct {
	Validation map[string]error
}

type UserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func NewUserClaims(user *entity.User) *UserClaims {
	return &UserClaims{
		ID:       user.ID,
		Username: user.Username,
	}
}
