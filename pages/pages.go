package pages

import (
	"errors"
	"fmt"
)

// pages file names as constant, just to keep things organized
const (
	Index         = "index.tmpl"
	Error         = "error.tmpl"
	Login         = "login.tmpl"
	WhatsappLogin = "whatsapp_login.tmpl"
)

// ==============================================================================
// template functions
// ==============================================================================

var TmplFunc map[string]any = map[string]any{
	"Map": Map,
}

func Map(args ...any) (map[string]any, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("map must have even number of arguments")
	}

	m := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		m[fmt.Sprintf("%v", args[i])] = args[i+1]
	}

	return m, nil
}
