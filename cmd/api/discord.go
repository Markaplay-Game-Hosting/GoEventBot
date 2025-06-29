package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type DiscordBody struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

type Embed struct {
	Color       int    `json:"color"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeStamps  string `json:"timestamp"`
}

func (app *application) SendMessage(embeds []Embed, title string, webhookId uuid.UUID) error {

	body := DiscordBody{
		Content: title,
		Embeds:  embeds,
	}
	bodyJson, err := json.Marshal(body)
	app.logger.Info("body: ", bodyJson)
	if err != nil {
		app.logger.Error("Unable to format body to send the message", err.Error())
		return err
	}
	webhook, err := app.models.Webhooks.GetByID(webhookId)
	if err != nil {
		app.logger.Error("Unable to get webhook by ID", "Error", err.Error())
		return fmt.Errorf("unable to get webhook by ID: %w", err)
	}
	resp, err := http.Post(webhook.URL, "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		app.logger.Error("Unable to send message", err.Error())
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			app.logger.Error("Unable to close body")
		}
	}(resp.Body)

	if resp.StatusCode != 204 {
		app.logger.Error("Unable to send message: ", resp.Status)
		return err
	}
	return nil
}

func FormatMessage(event data.Event) []Embed {
	var embed Embed
	var embeds []Embed
	rrule, err := ParseRRule(event.RRule)
	if err != nil {
		fmt.Println("Error parsing RRule:", err)
		return nil
	}
	embed.Title = event.Title
	embed.Description = event.Description
	// https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812
	embed.Color = 15105570
	embed.TimeStamps = rrule.After(time.Now(), false).Format(time.RFC3339)
	embeds = append(embeds, embed)
	return embeds
}
