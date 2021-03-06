package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"strconv"
)

func (bot *Bot) Give(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var id string
	if args["User"] == "" {
		id = message.Author.ID
	} else if len(args["User"]) == 18 {
		id = args["User"]
	} else if len(message.Mentions) > 0 {
		id = message.Mentions[0].ID
	} else {
		return UserNotFoundError
	}
	item, _ := strconv.Atoi(args["Item"])
	amount, _ := strconv.Atoi(args["Amount"])
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.Members.Apply(id, func(member *structures.Member) {
		member.Chests[uint(item)] += uint(amount)
	})
	session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
