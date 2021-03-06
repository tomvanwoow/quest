package experience

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	quest "github.com/tomvanwoow/quest/commands"
	"github.com/tomvanwoow/quest/utility"
	"time"
)

func GrantExp(bot *quest.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	s := struct {
		Guild  string
		Member string
	}{
		Guild:  utility.MustGetGuildID(session, message),
		Member: message.Author.ID,
	}
	if s.Guild == "" {
		return
	}
	t, ok := bot.ExpTimes[s]
	guild := bot.Guilds.Get(s.Guild)
	member := guild.Members.Get(s.Member)
	if !ok || uint16(time.Since(t).Seconds()) > guild.ExpReload {
		bot.ExpTimes[s] = time.Now()
		var r int64
		if guild.ExpGainLower > guild.ExpGainUpper {
			r = int64(utility.RandInt(int(guild.ExpGainUpper), int(guild.ExpGainLower)))
		} else {
			r = int64(utility.RandInt(int(guild.ExpGainLower), int(guild.ExpGainUpper)))
		}
		member.Experience += r
		fmt.Println(s.Member, r)
		GrantLottery(session, guild, member)
		GrantRoles(session, message, guild, member)
	}
}
