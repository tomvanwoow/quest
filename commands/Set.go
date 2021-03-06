package commands

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/structs"
	commands "github.com/tomvanwoow/disgone"
	"github.com/tomvanwoow/quest/modlog"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getOptions() map[string]*structures.Option {
	t := reflect.TypeOf(structures.Guild{})
	options := make(map[string]*structures.Option)
	for _, s := range structs.Names(structures.Guild{}) {
		field, _ := t.FieldByName(s)
		t, ok := field.Tag.Lookup("type")
		if ok {
			name, _ := field.Tag.Lookup("db")
			options[s] = &structures.Option{
				Name: name,
				Type: t,
			}
		}
	}
	return options
}

func (bot *Bot) Set(session *discordgo.Session, message *discordgo.MessageCreate, args map[string]string) error {
	options := getOptions()
	fmt.Println(options)
	if args["Option"] == "" {
		names := make([]string, len(options))
		var i int
		for name := range options {
			names[i] = name
			i++
		}
		sort.Strings(names)
		guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
		guild.RLock()
		var buf bytes.Buffer
		for _, name := range names {
			current := repr(reflect.Indirect(reflect.ValueOf(guild).Elem()).FieldByName(name).Interface(), options[name].Type)
			buf.WriteString(fmt.Sprintf("**%s** - %s\n", name, current))
		}
		guild.RUnlock()
		_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			Description: buf.String(),
		})
	} else if args["Value"] == "" {
		return commands.UsageError{
			Usage: bot.Commands["set"].GetUsage(bot.Prefix, "set"),
		}
	} else {
		keyName := args["Option"]
		option, ok := options[keyName]
		if !ok {
			var found bool
			for k, o := range options {
				if strings.ToLower(args["Option"]) == strings.ToLower(k) {
					option = o
					found = true
					keyName = k
					break
				}
			}
			if !found {
				return fmt.Errorf(`The provided argument for the Type was incorrect:
%s is **not** a Type.
Use q:types to view all types.`, args["Type"])
			}
		}
		pattern, _ := bot.Types[option.Type]
		value := args["Value"]
		result, _ := regexp.MatchString(pattern, value)
		if !result {
			return commands.UsageError{
				Usage: bot.Commands["set"].GetUsage(bot.Prefix, "set"),
			}
		}
		guild := bot.Guilds.Get(utility.MustGetGuildID(session, message))
		guild.RLock()
		defer guild.RUnlock()
		field := reflect.ValueOf(guild).Elem().FieldByName(keyName)
		val := reflect.ValueOf(convertType(message, option.Type, value)).Convert(field.Type())
		bot.Guilds.Apply(guild.ID, func(guild *structures.Guild) {
			field.Set(val)
		})
		if val.Type() == reflect.TypeOf(modlog.Modlog{}) {
			go guild.Modlog.StartLogging(session, &guild.Cases)
		}
		if guild.Modlog.Valid {
			guild.Modlog.Log <- &modlog.CaseSet{
				ModeratorID: message.Author.ID,
				Option:      keyName,
				Value:       repr(val.Interface(), option.Type),
			}
		}
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, "☑")
	}
	return nil
}

func repr(val interface{}, T string) string {
	switch v := val.(type) {
	case sql.NullString:
		if v.Valid {
			switch T {
			case "ChannelMention":
				return "<#" + v.String + ">"
			case "UserMention":
				return "<@" + v.String + ">"
			case "RoleMention":
				return "<@&" + v.String + ">"
			}
		} else {
			return "None"
		}
	case modlog.Modlog:
		if v.Valid {
			return "<#" + v.ChannelID + ">"
		} else {
			return "None"
		}
	}
	return fmt.Sprint(val)
}

func convertType(message *discordgo.MessageCreate, T string, value string) interface{} {
	var a interface{}
	switch T {
	case "Integer", "SignedInteger":
		a, _ = strconv.Atoi(value)
	case "Decimal", "Float", "BigNumber":
		c := strings.Split(value, "e")
		v, _ := strconv.ParseFloat(c[0], 32)
		if len(c) == 1 {
			a = v
		}
		e, _ := strconv.Atoi(c[1])
		a = v * math.Pow10(e)
	case "UserMention":
		if value == "none" {
			a = sql.NullString{}
		} else if len(message.Mentions) > 0 {
			a = sql.NullString{
				String: message.Mentions[0].ID,
				Valid:  true,
			}
		} else {
			a = sql.NullString{
				String: value,
				Valid:  true,
			}
		}
	case "RoleMention":
		if value == "none" {
			a = sql.NullString{}
		} else if len(message.MentionRoles) > 0 {
			a = sql.NullString{
				String: message.MentionRoles[0],
				Valid:  true,
			}
		} else {
			a = sql.NullString{
				String: value,
				Valid:  true,
			}
		}
	case "ChannelMention":
		if value == "none" {
			a = modlog.Modlog{}
		} else if len(value) > 18 {
			a = modlog.Modlog{
				ChannelID: value[2:20],
				Valid:     true,
				Log:       make(chan modlog.Case),
				Quit:      make(chan struct{}),
			}
		} else {
			a = modlog.Modlog{
				ChannelID: value,
				Valid:     true,
				Log:       make(chan modlog.Case),
				Quit:      make(chan struct{}),
			}
		}
	case "Boolean":
		l := strings.ToLower(value)
		if l == "true" || l == "yes" || l == "y" {
			a = true
		} else {
			a = false
		}
	}
	return a
}
