package models

import "github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"

// ListSettingsResponse
// @Description List of settings
type ListSettingsResponse struct {
	Settings []data.Setting `json:"settings"`
} // @name Settings.List

// AddSettingRequest
// @Description Request body to add a setting
type AddSettingRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	IsSecret    bool   `json:"is_secret"`
} // @name Settings.Add.Request

// UpdateSettingRequest
// @Description Request body to update a setting
type UpdateSettingRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
	IsSecret    bool   `json:"is_secret,omitempty"`
} // @name Settings.Update.Request
