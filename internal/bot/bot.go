package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/crypter"
	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/data"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	Logger  *slog.Logger
	Models  data.Models
	Crypter *crypter.Crypt
}

func New(logger *slog.Logger, models data.Models, crypt *crypter.Crypt) (*Bot, error) {

	return &Bot{
		Logger:  logger,
		Models:  models,
		Crypter: crypt,
	}, nil
}

func (b *Bot) Run(ctx context.Context) {

	var token string
	goto WAIT
WAIT:
	for {
		select {
		case <-ctx.Done():
			b.Logger.Info("discord bot is shutting down")
			return
		default:
			t, err := b.Models.Settings.GetByName("discord_token")
			if err == nil && t.Value != "" {
				b.Logger.Info("token detected!")
				tokenDecrypted, err := b.Crypter.Decrypt(t.Value)
				if err != nil {
					b.Logger.Error("Unable to decrypt token", "error", err)
				}
				if token != tokenDecrypted {
					token = tokenDecrypted
				}

				goto READY
			}
			b.Logger.Info("Checking for token in 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
READY:
	session, err := discordgo.New("Bot " + token)
	b.Session = session
	if err != nil {
		b.Logger.Error("Error creating Discord session: ", "error", err)
		goto WAIT
	}
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.Logger.Info(fmt.Sprintf("Logged in as %s", r.User.String()))
	})

	if err := b.Session.Open(); err != nil {
		b.Logger.Error("Error opening Discord session", "err", err)
		return
	}
	//defer b.Session.Close()

	<-ctx.Done()
	b.Logger.Info("Bot shutting down")
	b.Session.Close()
}

func (b *Bot) SendMessage(channelID string, embed *discordgo.MessageEmbed) error {
	_, err := b.Session.ChannelMessageSendEmbed(channelID, embed)

	return err
}

func (b *Bot) GetChannel(channelID string) (*discordgo.Channel, error) {
	return b.Session.Channel(channelID)
}

func (b *Bot) ListChannels(guildID string) ([]*discordgo.Channel, error) {
	return b.Session.GuildChannels(guildID)
}

func (b *Bot) ListRoles(guildId string) ([]*discordgo.Role, error) {
	return b.Session.GuildRoles(guildId)
}

func (b *Bot) GetGuildInfo(guildID string) (*discordgo.Guild, error) {
	return b.Session.Guild(guildID)
}

func FormatMessage(event data.EventInstance) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       event.Title,
		Description: event.Description,
		Color:       15105570,
		Timestamp:   event.StartDate.Format(time.RFC3339),
	}
}
