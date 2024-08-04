package main

import (
	"bytes"
	"encoding/json"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
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
	TimeStamps  string `json:"timeStamps"`
}

func (app *application) SendMessage(embeds []Embed) error {

	body := DiscordBody{
		Content: app.config.recipient,
		Embeds:  embeds,
	}
	bodyJson, err := json.Marshal(body)
	app.logger.Info("body: ", bodyJson)
	if err != nil {
		app.logger.Error("Unable to format body to send the message", err.Error())
		return err
	}
	resp, err := http.Post(app.config.webhook, "application/json", bytes.NewBuffer(bodyJson))
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

	embed.Title = event.Title
	embed.Description = event.Description
	// https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812
	embed.Color = 15105570
	embed.TimeStamps = event.StartDate.Format(time.RFC3339)
	embeds = append(embeds, embed)
	return embeds
}
