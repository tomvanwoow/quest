package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	commands "github.com/tomvanwoow/disgone"
	"github.com/tomvanwoow/quest/inventory"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/structures"
)

type Bot struct {
	*commands.BotOptions
	ExpTimes map[struct {
		Guild  string
		Member string
	}]time.Time
	structures.Guilds
	DB *sqlx.DB
	inventory.Chests
	Modlogs    map[string]modlog.Modlog
	ErrorEmbed func(e error) *discordgo.MessageEmbed
	Embed      func(title string,
		description string,
		fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed
}
