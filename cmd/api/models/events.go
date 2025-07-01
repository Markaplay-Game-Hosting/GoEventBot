package models

// CreateEventRequest
// @Description Information required to create a new event
type CreateEventRequest struct {
	// Title of the  event
	Title string `json:"title"`
	// Description/Content of the event
	Description string `json:"description"`
	// Duration of the each event following ISO 8601 standard
	Duration string `json:"duration" example:"PT30M"`
	// Recurrence rule following RFC 5545 https://icalendar.org/rrule-tool.html
	RRule string `json:"rrule" example:"FREQ=WEEKLY;INTERVAL=1;BYDAY=MO;UNTIL=20250731T000000Z"`
	// Discord ID of the channel the bot will post the event
	ChannelID int `json:"channel_id,omitempty"`
	// Discord Guild ID/server the bot will publish on
	GuildID int `json:"guild_id,omitempty"`
} // @name Event.Create.Request

// UpdateEventRequest
// @Description Information required to update a new event
type UpdateEventRequest struct {
	// Title of the  event
	Title string `json:"title,omitempty"`
	// Description/Content of the event
	Description string `json:"description,omitempty"`
	// Duration of the each event following ISO 8601 standard
	Duration string `json:"duration,omitempty" example:"PT30M"`
	// Recurrence rule following RFC 5545 https://icalendar.org/rrule-tool.html
	RRule string `json:"rrule,omitempty" example:"FREQ=WEEKLY;INTERVAL=1;BYDAY=MO;UNTIL=20250731T000000Z"`
	// Discord ID of the channel the bot will post the event
	ChannelID int `json:"channel_id,omitempty"`
	// Discord Guild ID/server the bot will publish on
	GuildID int `json:"guild_id,omitempty"`
	// If the event is active of not
	IsActive bool `json:"is_active,omitempty"`
} // @name Event.Create.Request
