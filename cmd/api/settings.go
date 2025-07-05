package main

import (
	"net/http"

	"github.com/Markaplay-Game-Hosting/GoEventBot/cmd/api/models"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
)

// listSettings
// @Summary      List Settings
// @Description  List all settings
// @Tags         Settings
// @Produce      json
// @Success      200 {object} models.ListSettingsResponse
// @Router       /jobs [get]
func (app *application) listSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := app.models.Settings.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"settings": settings}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getSetting
// @Summary      Get Setting
// @Description  get setting
// @Tags         Settings
// @Produce      json
// @Success      200 {object} data.Setting
// @Router       /jobs/{id} [get]
func (app *application) getSettings(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	settings, err := app.models.Settings.GetByID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"settings": settings}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// addSetting
// @Summary      Add Setting
// @Description  Add new settings
// @Tags         Settings
// @Param        setting body models.AddSettingRequest true "body to add a setting"
// @Produce      json
// @Success      201 {object} data.Setting
// @Router       /jobs [post]
func (app *application) addSettings(w http.ResponseWriter, r *http.Request) {
	var input models.AddSettingRequest

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Error("Unable to read JSON", "error", err)
		app.badRequestResponse(w, r, err)
	}

	if input.IsSecret {
		input.Value = app.crypt.Encrypt(input.Value)
	}

	setting := &data.Setting{
		Name:        input.Name,
		Description: input.Description,
		Value:       input.Value,
	}
	err = app.models.Settings.Insert(setting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"setting": setting}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateSetting
// @Summary      Update Setting
// @Description  Update existing setting
// @Tags         Settings
// @param        id path string true "id of the setting to update"
// @Param        setting body models.UpdateSettingRequest true "body to update a setting"
// @Produce      json
// @Success      200 {object} data.Setting
// @Router       /jobs/{id} [put]
func (app *application) updateSetting(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.logger.Error("Unable to read JSON", err.Error())
		app.badRequestResponse(w, r, err)
	}

	var input models.AddSettingRequest

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Error("Unable to read JSON", err.Error())
		app.badRequestResponse(w, r, err)
	}

	if input.IsSecret {
		input.Value = app.crypt.Encrypt(input.Value)
	}

	setting, err := app.models.Settings.GetByID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if input.Name != "" {
		setting.Name = input.Name
	}
	if input.Description != "" {
		setting.Description = input.Description
	}

	if input.Value != "" {
		setting.Value = input.Value
	}

	err = app.models.Settings.Update(setting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusNoContent, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteSettings
// @Summary      Delete Setting
// @Description  Delete an existing setting
// @Tags         Settings
// @param        id path string true  "id of the setting to delete"
// @Produce      json
// @Success      204 {object} data.Setting
// @Router       /jobs/{id} [delete]
func (app *application) deleteSetting(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.logger.Error("Unable to read JSON", err.Error())
		app.badRequestResponse(w, r, err)
	}
	err = app.models.Settings.Delete(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusNoContent, nil, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
