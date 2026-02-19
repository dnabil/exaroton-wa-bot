package pages

import (
	"github.com/CloudyKit/jet/v6"
)

// pages file names as constant, just to keep things organized
const (
	Index         = "index.jet"
	Error         = "error.jet"
	Login         = "login.jet"
	WhatsappLogin = "whatsapp_login.jet"

	SettingsExaroton = "settings_exaroton.jet"
	SettingsWhatsapp = "settings_whatsapp.jet"
)

// ==============================================================================
// template functions
// ==============================================================================

var TmplFunc map[string]jet.Func = map[string]jet.Func{}
