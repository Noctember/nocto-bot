package main

import (
	"Noctobot/commands"
	"Noctobot/database"
	"Noctobot/events"
	"Noctobot/utils"
	"bytes"
	"flag"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"github.com/spf13/cast"
	"time"
)

var (
	token  = flag.String("t", "", "Bot account token")
	initdb = flag.Bool("init", false, "Initializes the database.")
)

func main() {
	flag.Parse()

	if *initdb {
		utils.Logger.Info("Creating database schemas...")
		before := time.Now()
		database.MustExec(database.SCHEMA)
		after := time.Now()
		utils.Logger.Infof("Schema creation took %d ms", after.Sub(before).Milliseconds())
		return
	}

	token := utils.GetConfig("devtoken")

	if cast.ToBool(utils.GetConfig("dev")) {
		token = utils.GetConfig("token")
	}

	dg, err := discordgo.New("Bot " + token)
	/*	dg.GatewayManager = &impl.GatewayManager{
		session:          s,
		voiceConnections: make(map[int64]*VoiceConnection),
	}*/
	dg.State.TrackVoice = false
	panicOnErr(err)

	bot := gocto.New(dg)
	bot.LoadBuiltins()
	bot.AddMonitor(gocto.NewMonitor("qrcodeDecoder", events.QRCode))

	bot.OwnerID = 118499178790780929
	bot.Color = 0xDFAC7C
	bot.SetErrorHandler(utils.ErrorHandler)
	bot.SetListHandler(func(b *gocto.Bot, m *discordgo.Message) bool {
		if database.IsPlonked(m.Author.ID) {
			return true
		}
		return false
	})
	bot.SetMentionPrefix(true)
	bot.SetPrefixHandler(func(_ *gocto.Bot, msg *discordgo.Message, dm bool) string {
		if bytes.HasPrefix([]byte(msg.Content), []byte("nocto ")) {
			return "nocto "
		} else if dm {
			return "="
		}
		return database.GetPrefix(msg.GuildID)
	})

	bot.SetInvitePerms(268552262)

	commands.Init(bot)
	events.Init(dg)
	events.Listen(bot)

	bot.MustConnect()

	bot.Wait()

	defer database.Close()
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
