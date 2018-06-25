package events

import (
	"github.com/bwmarrin/discordgo"
	commands "../../discordcommands"
)

func MemberAdd(bot *commands.Bot) func(*discordgo.Session, *discordgo.GuildMemberAdd) {
	return func(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
		guild := bot.Guilds.Get(event.GuildID)
		if guild.Autorole.Valid {
			session.GuildMemberRoleAdd(event.GuildID, event.Member.User.ID, guild.Autorole.String)
		}
	}
}