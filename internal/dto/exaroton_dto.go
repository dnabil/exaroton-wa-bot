package dto

type ExarotonAccountInfo struct {
	// Name represents the account's name.
	Name string `json:"name"`

	// Email represents the account's email.
	Email string `json:"email"`

	// Verified represents whether the account is verified.
	Verified bool `json:"verified"`

	// Credits represents the account's credits.
	Credits float64 `json:"credits"`
}
