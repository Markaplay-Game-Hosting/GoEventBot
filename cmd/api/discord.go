package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/cmd/api/models"
	"github.com/bwmarrin/discordgo"
	"net/http"
)

// listRolesHandler
// @Summary      List Roles from a guild
// @Description  List Roles from a guild
// @Tags         Discord
// @Produce      json
// @Success      200 {array} models.Role
// @Router       /discord/{guild}/roles [get]
func (app *application) listRolesHandler(w http.ResponseWriter, r *http.Request) {
	guildId := app.readString(r, "id", "")
	if guildId == "" {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
	}
	app.logger.Info("guildId", "ID", guildId)
	allRoles, err := app.bot.ListRoles(guildId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	var roles []models.Role
	for _, role := range allRoles {
		mRole := models.Role{
			ID:          role.ID,
			Name:        role.Name,
			Mentionable: role.Mentionable,
		}
		roles = append(roles, mRole)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"roles": roles}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listChannelsHandler
// @Summary      List Channels from a guild
// @Description  List Channels from a guild
// @Tags         Discord
// @Produce      json
// @Success      200 {array} models.Channel
// @Router       /discord/{guild}/channels [get]
func (app *application) listChannelsHandler(w http.ResponseWriter, r *http.Request) {
	guildId := app.readString(r, "id", "")
	if guildId == "" {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
	}
	app.logger.Info("guildId", "ID", guildId)
	allChannels, err := app.bot.ListChannels(guildId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	var channels []models.Channel
	for _, channel := range allChannels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		mChannel := models.Channel{
			ID:   channel.ID,
			Name: channel.Name,
		}
		channels = append(channels, mChannel)
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"channels": channels}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getGuildInfo
// @Summary      Get Guild Info
// @Description  Show guild basic information
// @Tags         Discord
// @Produce      json
// @Success      200 {array} models.Guild
// @Router       /discord/{guild} [get]
func (app *application) getGuildInfo(w http.ResponseWriter, r *http.Request) {
	guildId := app.readString(r, "id", "")
	if guildId == "" {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
	}
	app.logger.Info("guildId", "ID", guildId)
	guildInfo, err := app.bot.GetGuildInfo(guildId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	guild := models.GuildInfo{
		ID:   guildId,
		Name: guildInfo.Name,
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"guild": guild}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
