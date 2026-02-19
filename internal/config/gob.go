package config

import (
	"encoding/gob"
	"exaroton-wa-bot/internal/dto"
)

var registeredGobs []any = []any{
	// add types to be registered here
	dto.WebFlashMessage{},
	dto.WebValidationErrors{},
	dto.WebOldInput{},
}

func RegisterGobs() {
	for _, t := range registeredGobs {
		gob.Register(t)
	}
}
