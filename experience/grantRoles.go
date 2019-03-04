package experience

import (
<<<<<<< HEAD:src/quest/experience/grantRoles.go
	"../structures"
=======
	commands "github.com/tomvanwoow/discordcommands"
	"github.com/tomvanwoow/quest/structures"
>>>>>>> master:experience/grantRoles.go
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
	"../utility"
)

func GrantRoles(session *discordgo.Session, message *discordgo.MessageCreate, guild *structures.Guild, member *structures.Member) error {
	m, err := session.GuildMember(guild.ID, member.ID)
	if err != nil {
		return err
	}
	allRoles, err := session.GuildRoles(guild.ID)
	if err != nil {
		return nil
	}
	discordRoles := make(discordgo.Roles, len(guild.Roles))
	for i, r := range guild.Roles {
		ok, index := roleContains(allRoles, r.ID)
		if ok {
			discordRoles[i] = allRoles[index]
		}
	}
	for i, role := range discordRoles {
		r := guild.Roles[i]
		if member.Experience >= r.Experience {
			if role == nil {
				continue
			}
			found, _ := utility.Contains(m.Roles, role.ID)
			if !found {
				session.GuildMemberRoleAdd(guild.ID, member.ID, r.ID)
				session.ChannelMessageSendEmbed(message.ChannelID, questEmbedColor(m.User.Username, m.User.Discriminator, role.Name, role.Color))
			}
		}
	}
	return nil
}

func questEmbedColor(username string, discriminator string, rolename string, color int) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("Congratulations %s#%s", username, discriminator),
		Description: fmt.Sprintf("You received the %s role", rolename),
		Timestamp:   utility.TimeToTimestamp(time.Now()),
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
	}
	return emb
}

func roleContains(roles []*discordgo.Role, id string) (bool, int) {
	for i, r := range roles {
		if r.ID == id {
			return true, i
		}
	}
	return false, 0
}