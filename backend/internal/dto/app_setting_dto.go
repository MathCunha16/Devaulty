package dto

type AppSettingView struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateAppSettingCommand struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
