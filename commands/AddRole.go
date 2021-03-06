package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"strconv"
)

func (bot *Bot) AddRole(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	var roleID string
	if len(message.MentionRoles) > 0 {
		roleID = message.MentionRoles[0]
	} else {
		roleID = args["Role"]
	}
	guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
	guild.RLock()
	defer guild.RUnlock()
	exp, _ := strconv.Atoi(args["Experience"])
	if exp == 0 {
		return fmt.Errorf(`I see you are trying to add a role for 0 experience!
If you want to, it is better to use q:set Autorole <role>`)
	}
	allIDs := make([]string, len(guild.Roles))
	for i, v := range guild.Roles {
		allIDs[i] = v.ID
	}
	found, index := utility.Contains(allIDs, roleID)
	if !found && len(guild.Roles) >= 64 {
		return fmt.Errorf("Invalid action - 64 roles is the absolute limit\nTry removing a role")
	}
	role := &structures.Role{
		Experience: int64(exp),
		ID:         roleID,
	}
	bot.Guilds.Apply(guild.ID, func (guild *structures.Guild) {
		if found {
			guild.Roles[index] = role
		} else {
			guild.Roles = append(guild.Roles, role)
		}
	})
	_ = session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	return nil
}
