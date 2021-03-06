package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/modlog"
	"time"
)

func (bot BotEvents) GuildCreate(session *discordgo.Session, guild *discordgo.GuildCreate) {
	session.State.GuildAdd(guild.Guild)
	time.Sleep(time.Second)
	g := bot.Guilds[guild.ID]
	if g != nil && g.Modlog.Valid {
		go modlog.StartLogging(session, g.Modlog, &g.Cases)
	}
}
