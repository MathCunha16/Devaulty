package dto

type MasterPassword struct {
	MasterPassword string `json:"masterPassword" binding:"required,min=8"`
}

type VaultStatusView struct {
	Active      bool  `json:"active"`
	SecondsLeft int64 `json:"secondsLeft"`
}
