package bot

import (
	"os"

	"github.com/Markaplay-Game-Hosting/GoEventBot/internal/config"
	"github.com/bwmarrin/discordgo"
)

func Run(cfg config.Config) {
	session, err := discordgo.New("Bot" + cfg.Discord.ClientID)
	if err != nil {
		panic("Error creating Discord session: " + err.Error())
		os.Exit(1)
	}
	session.AddHandler(SendMessage)
}
