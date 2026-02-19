package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ServerSettings struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingsExarotonPageData struct {
	HttpCode   int
	Validation map[string]error

	APIKey      string
	APIKeyMsg   string
	AccountInfo *ExarotonAccountInfo
}

type SettingsExarotonReq struct {
	APIKey string `json:"api_key"`
}

func (r *SettingsExarotonReq) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.APIKey, validation.Required, validation.Length(10, 200)),
	)
}
