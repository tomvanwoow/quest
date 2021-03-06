package commands

import (
	commands "github.com/tomvanwoow/disgone"
	"github.com/bwmarrin/discordgo"
	"github.com/tomvanwoow/quest/utility"
)

const (
	PermissionMember commands.Group = iota
	PermissionModerator
	PermissionAdmin
	PermissionOwner
)

func (bot *Bot) UserGroup(session *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member) commands.Group {
	g := bot.Guilds.Get(guild.ID)
	if member.User.ID == guild.OwnerID {
		return PermissionOwner
	}
	for _, r := range member.Roles {
		role, err := utility.GetRole(session, guild.ID, r)
		if err == nil {
			if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
				return PermissionAdmin
			}
		}
	}
	for _, r := range member.Roles {
		if r == g.AdminRole.String {
			return PermissionAdmin
		} else if r == g.ModRole.String {
			return PermissionModerator
		}
	}
	return PermissionMember
}
