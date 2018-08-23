package commands

import (
	commands "../../discordcommands"
	"../modlog"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) Unban(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var user *discordgo.User
	if len(args["User"]) == 18 {
		var err error
		user, err = session.User(args["User"])
		if err != nil {
			return UserNotFoundError{}
		}
	} else if len(message.Mentions) > 0 {
		user = message.Mentions[0]
	} else {
		return UserNotFoundError{}
	}
	err := session.GuildBanDelete(commands.MustGetGuildID(session, message), user.ID)
	if err != nil {
		return fmt.Errorf("Can't ban that user! Make sure I have the discord ban permission.")
	}
	guild := bot.Guilds.Get(commands.MustGetGuildID(session, message))
	if guild.Modlog.Valid {
		guild.Modlog.Log <- &modlog.CaseUnban{
			ModeratorID: message.Author.ID,
			UserID:      user.ID,
			Reason:      args["Reason"],
		}
	}
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
