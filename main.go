package main

import (
	_ "database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	commands "github.com/tomvanwoow/disgone"
	quest "github.com/tomvanwoow/quest/commands"
	"github.com/tomvanwoow/quest/events"
	"github.com/tomvanwoow/quest/inventory"
	"github.com/tomvanwoow/quest/structures"
	"github.com/tomvanwoow/quest/utility"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

var CommandsData commands.CommandMap
var Types = map[string]string{}

type App struct {
	Token    string
	User     string
	Pass     string
	Host     string
	Database string
	Commands string
	Types    string
}

var db *sqlx.DB
var app App
var chests inventory.Chests
var bot *quest.Bot

const (
	prefix = "q:"
)

func questEmbed(title string, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:      "rich",
		Title:     title,
		Timestamp: utility.TimeToTimestamp(time.Now().UTC()),
		Color:     0x00ffff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
		Fields: fields,
	}
	if description != "" {
		emb.Description = description
	}
	return emb
}

func errorEmbed(e error) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "An error has occurred",
		Description: e.Error(),
		Timestamp:   utility.TimeToTimestamp(time.Now()),
		Color:       0x660000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Quest Bot",
		},
	}
	return emb
}

func unmarshalJson(filename string, v interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return
	}
	data := make([]byte, stat.Size())
	_, err = f.Read(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	var err error
	var src string
	flag.StringVar(&src, "a", "", "App Location")
	flag.Parse()
	err = unmarshalJson(src, &app)
	if err != nil {
		log.Fatalln(err)
	}
	err = unmarshalJson(app.Commands, &CommandsData)
	if err != nil {
		log.Fatalln(err)
	}
	err = unmarshalJson(app.Types, &Types)
	if err != nil {
		log.Fatalln(err)
	}
	err = unmarshalJson("src/json/chests.json", &chests)
	if err != nil {
		log.Fatalln(err)
	}
	db, err = structures.InitDB(app.User, app.Pass, app.Host, app.Database)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	defer db.Close()
	bot = &quest.Bot{
		ExpTimes: make(map[struct {
			Guild  string
			Member string
		}]time.Time),
		Guilds:     structures.NewGuildsCache(db),
		DB:         db,
		Chests:     chests,
		Embed:      questEmbed,
		ErrorEmbed: errorEmbed,
		BotOptions: &commands.BotOptions{
			Commands: CommandsData,
			Errors: make(chan struct {
				Err error
				*discordgo.MessageCreate
			}),
			Types: Types,
			GroupNames: map[commands.Group]string{
				quest.PermissionMember:    "Member",
				quest.PermissionModerator: "Moderator",
				quest.PermissionAdmin:     "Admin",
				quest.PermissionOwner:     "Owner",
			},
			Prefix: prefix,
			OnPanic: func(session *discordgo.Session, message *discordgo.MessageCreate, r interface{}) {
				log.Println(string(debug.Stack()))
				session.ChannelMessageSend(message.ChannelID, "```"+`An unexpected panic occured in the execution of that command.
`+fmt.Sprint(r)+"\nTry again later, or contact themeeman#8354"+"```")
			},
		},
	}
	dg, err := commands.NewSession(bot, bot.BotOptions, app.Token)
	if err != nil {
		log.Fatalln("Error making discord session", err)
	}
	e := events.BotEvents{Bot: bot}
	dg.AddHandlerOnce(e.Ready)
	dg.AddHandler(func(session *discordgo.Session, _ *discordgo.Ready) { session.UpdateStatus(0, "q:help") })
	dg.AddHandler(e.MessageCreate)
	dg.AddHandler(e.MemberAdd)
	dg.AddHandler(e.GuildCreate)
	dg.StateEnabled = true
	dg.State.TrackMembers = true
	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection", err)
	}
	defer dg.Close()
	fmt.Println("Quest is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	structures.CommitAllGuilds(bot.Guilds)
}
