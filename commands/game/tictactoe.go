package game

import (
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ONE   = ":one:"
	TWO   = ":two:"
	THREE = ":three:"
	FOUR  = ":four:"
	FIVE  = ":five:"
	SIX   = ":six:"
	SEVEN = ":seven:"
	EIGHT = ":eight:"
	NINE  = ":nine:"
)

var ireg = regexp.MustCompile(`^([1-9]|one|two|three|four|five|six|seven|eigth|nine)\s*`)

type tGame struct {
	Grid      tGrid
	ChannelID int64
	Turn      int
	Stop      chan bool
	Win       bool
	Started   bool
	Players   []int64
	MessageID int64
	Coin      map[int64]string
}

type tGrid struct {
	Line1   []int
	Line2   []int
	Line3   []int
	Message string
}

type activeGames []tGame

func TicTacToe(ctx *gocto.CommandContext) {
	game := tGame{
		Grid:      tGrid{},
		Coin:      make(map[int64]string),
		ChannelID: ctx.Channel.ID,
		Turn:      0,
		Win:       false,
		Stop:      make(chan bool),
		Started:   false,
		Players:   []int64{ctx.Author.ID},
	}
	embed := gocto.NewEmbed().
		SetDescription(utils.CheckMark+" A new game of tic tac toe has been started!\nStarting now, you have 5 minutes to start and finish the game.").
		SetFooter("React with 🚪 to join").
		AddField("Players", ctx.Author.Username+"\n")
	msg, _ := ctx.ReplyEmbed(embed.Build())
	game.MessageID = msg.ID
	ctx.Session.MessageReactionAdd(ctx.Channel.ID, msg.ID, ":door:")
	go game.timeOut(ctx)
	if waitForPlayers(ctx, &game) {
		waitForPlay(ctx, &game)
	}
}

type play struct {
	UserID int64
	Move   string
}

func waitForPlay(ctx *gocto.CommandContext, game *tGame) {
	for {
		var r *play
		select {
		case <-game.Stop:
			return
		case e := <-nextPlay(ctx.Bot, game):
			r = e
		}
		if !contains(game.Players, r.UserID) {
			continue
		}
		i := 0
		nexti := 0
		for in, id := range game.Players {
			if id == r.UserID {
				i = in + 1
				nexti = i + 1
			}
		}
		if game.Turn != i {
			continue
		}
		if len(game.Players) < nexti {
			nexti = 1
		}

		r.Move = strings.ToLower(r.Move)
		r.Move = strings.ReplaceAll(r.Move, " ", "")
		r.Move = strings.ReplaceAll(r.Move, "1️⃣", "1")
		r.Move = strings.ReplaceAll(r.Move, "2️⃣", "2")
		r.Move = strings.ReplaceAll(r.Move, "3️⃣", "3")
		r.Move = strings.ReplaceAll(r.Move, "4️⃣", "4")
		r.Move = strings.ReplaceAll(r.Move, "5️⃣", "5")
		r.Move = strings.ReplaceAll(r.Move, "6️⃣", "6")
		r.Move = strings.ReplaceAll(r.Move, "7️⃣", "7")
		r.Move = strings.ReplaceAll(r.Move, "8️⃣", "8")
		r.Move = strings.ReplaceAll(r.Move, "9️⃣", "9")
		if !ireg.MatchString(r.Move) {
			continue
		}

		r.Move = strings.ReplaceAll(r.Move, "one", "1")
		r.Move = strings.ReplaceAll(r.Move, "two", "2")
		r.Move = strings.ReplaceAll(r.Move, "three", "3")
		r.Move = strings.ReplaceAll(r.Move, "four", "4")
		r.Move = strings.ReplaceAll(r.Move, "five", "5")
		r.Move = strings.ReplaceAll(r.Move, "six", "6")
		r.Move = strings.ReplaceAll(r.Move, "seven", "7")
		r.Move = strings.ReplaceAll(r.Move, "eight", "8")
		r.Move = strings.ReplaceAll(r.Move, "nine", "9")

		NumString := []string{"", ":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":nine:"}
		number, _ := strconv.Atoi(r.Move)
		if number > 0 && number < 10 {
			UpdateGrid := false
			if number > 6 {
				number -= 7
				if game.Grid.Line1[number] == 0 {
					UpdateGrid = true
					game.Grid.Line1[number] = i
				}
			} else if number > 3 {
				number -= 4
				if game.Grid.Line2[number] == 0 {
					UpdateGrid = true
					game.Grid.Line2[number] = i
				}
			} else {
				number -= 1
				if game.Grid.Line3[number] == 0 {
					UpdateGrid = true
					game.Grid.Line3[number] = i
				}
			}
			if UpdateGrid {
				Win, Tie := game.checkForWin(i)
				number, _ := strconv.Atoi(r.Move)
				user := ctx.User(r.UserID)
				NewContent := strings.Replace(game.Grid.Message, NumString[number], game.Coin[r.UserID], 1)

				game.Grid.Message = NewContent
				game.Turn = nexti

				next := ctx.Member(game.Players[nexti-1])
				if Win {
					ctx.Reply(NewContent)
					ctx.ReplyNoEdit(":crown: %s won the game!", user.Username)
					return
				} else if Tie {
					ctx.Reply(NewContent)
					ctx.ReplyNoEdit("No one won the game!")
					return
				} else {
					ctx.Reply("It's %s's turn!\n\n"+NewContent, next.User.Mention())
				}
			}

		}
	}
}

func (game *tGame) checkForWin(number int) (bool, bool) {
	if game.Grid.Line1[0] == number && game.Grid.Line1[1] == number && game.Grid.Line1[2] == number {
		return true, false
	}
	if game.Grid.Line2[0] == number && game.Grid.Line2[1] == number && game.Grid.Line2[2] == number {
		return true, false
	}
	if game.Grid.Line3[0] == number && game.Grid.Line3[1] == number && game.Grid.Line3[2] == number {
		return true, false
	}

	if game.Grid.Line1[0] == number && game.Grid.Line2[0] == number && game.Grid.Line3[0] == number {
		return true, false
	}
	if game.Grid.Line1[1] == number && game.Grid.Line2[1] == number && game.Grid.Line3[1] == number {
		return true, false
	}
	if game.Grid.Line1[2] == number && game.Grid.Line2[2] == number && game.Grid.Line3[2] == number {
		return true, false
	}

	if game.Grid.Line1[0] == number && game.Grid.Line2[1] == number && game.Grid.Line3[2] == number {
		return true, false
	}
	if game.Grid.Line1[2] == number && game.Grid.Line2[1] == number && game.Grid.Line3[0] == number {
		return true, false
	}
	if game.checkForTie() {
		return false, true
	}
	return false, false
}

func nextPlay(bot *gocto.Bot, game *tGame) chan *play {
	channel := make(chan *play)
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.ChannelID == game.ChannelID && r.MessageID == game.MessageID {
			channel <- &play{
				UserID: r.UserID,
				Move:   r.Emoji.Name,
			}
		}
	})
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageCreate) {
		if r.ChannelID == game.ChannelID {
			channel <- &play{
				UserID: r.Author.ID,
				Move:   r.Message.Content,
			}
		}
	})
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionRemove) {
		if r.ChannelID == game.ChannelID && r.MessageID == game.MessageID {
			channel <- &play{
				UserID: r.UserID,
				Move:   r.Emoji.Name,
			}
		}
	})
	return channel
}

func waitForPlayers(ctx *gocto.CommandContext, game *tGame) bool {
	for {
		var r *discordgo.MessageReaction

		select {
		case <-game.Stop:
			return false
		case e := <-nextReaction(ctx.Bot, game):
			r = e
		}

		user := ctx.Member(r.UserID)
		if user.User.Bot {
			continue
		}

		switch e := r.Emoji; e.Name {
		case "🚪":
			if game.Players[0] == user.User.ID {
				continue
			}

			if contains(game.Players, user.User.ID) {
				game.removePlayer(ctx, user.User)
			} else {
				game.addPlayer(ctx, user.User)
			}

			game.checkStart(ctx)
		case "⭐":
			if !contains(game.Players, r.UserID) {
				continue
			}
			if game.checkStart(ctx) {
				game.Grid = tGrid{
					Line1: []int{0, 0, 0},
					Line2: []int{0, 0, 0},
					Line3: []int{0, 0, 0},
					Message: ":one: :two: :three:\n" +
						":four: :five: :six:\n" +
						":seven: :eight: :nine:\n",
				}
				ctx.Session.ChannelMessageDelete(game.ChannelID, game.MessageID)
				m := map[int]string{
					0:  ":apple:",
					1:  ":pear:",
					2:  ":avocado:",
					3:  ":tangerine:",
					4:  ":lemon:",
					5:  ":banana:",
					6:  ":watermelon:",
					7:  ":grapes:",
					8:  ":strawberry:",
					9:  ":strawberry:",
					10: ":melon:",
					11: ":cherries:",
					12: ":peach:",
					13: ":mango:",
					14: ":pineapple:",
					15: ":coconut:",
					16: ":kiwi:",
					17: ":tomato:",
					18: ":eggplant:",
				}
				for i, u := range game.Players {
					game.Coin[u] = m[i]
					if u == r.UserID {
						game.Turn = i + 1
					}
				}
				msg, _ := ctx.Session.ChannelMessageSend(game.ChannelID, "It's "+user.User.Mention()+"'s turn!\n\n"+game.Grid.Message)
				ctx.Bot.CommandEdits[ctx.Message.ID] = msg.ID
				game.MessageID = msg.ID

				go func() {
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, ONE)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, TWO)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, THREE)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, FOUR)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, FIVE)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, SIX)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, SEVEN)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, EIGHT)
					ctx.Session.MessageReactionAdd(game.ChannelID, msg.ID, NINE)
				}()
				return true
			} else {
				ctx.Reply("%s Not enough players to start.", utils.CrossMark)
			}
		}
	}
}

func (g *tGame) addPlayer(ctx *gocto.CommandContext, user *discordgo.User) {
	g.Players = append(g.Players, user.ID)
	msg, _ := ctx.Session.ChannelMessage(g.ChannelID, g.MessageID)
	embed := msg.Embeds[0]
	embed.Fields[0].Value = embed.Fields[0].Value + "\n" + user.Username
	ctx.ReplyEmbed(embed)
}

func (g *tGame) removePlayer(ctx *gocto.CommandContext, user *discordgo.User) {
	for i, p := range g.Players {
		if p == user.ID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			msg, _ := ctx.Session.ChannelMessage(g.ChannelID, g.MessageID)
			embed := msg.Embeds[0]
			embed.Fields[0].Value = strings.ReplaceAll(embed.Fields[0].Value, "\n"+user.Username, "")
			ctx.ReplyEmbed(embed)

		}
	}
}

func (g *tGame) checkStart(ctx *gocto.CommandContext) bool {
	if len(g.Players) > 1 {
		ctx.Session.MessageReactionAdd(g.ChannelID, g.MessageID, "⭐")
		return true
	} else {
		ctx.Session.MessageReactionRemoveMe(g.ChannelID, g.MessageID, "⭐")
		return false
	}
}

func (game *tGame) checkForTie() bool {
	for i := 0; i <= 2; i++ {
		if game.Grid.Line1[i] == 0 {
			return false
		}
		if game.Grid.Line2[i] == 0 {
			return false
		}
		if game.Grid.Line3[i] == 0 {
			return false
		}
	}
	return true
}

func (game *tGame) timeOut(ctx *gocto.CommandContext) {
	time.Sleep(5 * time.Minute)
	if !game.Win {
		ctx.Reply("Too slow, the game has been deleted.")
		game.Stop <- true
	}
	return
}

func nextReaction(bot *gocto.Bot, game *tGame) chan *discordgo.MessageReaction {
	channel := make(chan *discordgo.MessageReaction)
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.ChannelID == game.ChannelID && r.MessageID == game.MessageID {
			channel <- r.MessageReaction
		}
	})
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionRemove) {
		if r.ChannelID == game.ChannelID && r.MessageID == game.MessageID {
			channel <- r.MessageReaction
		}
	})
	return channel
}

func contains(s []int64, u int64) bool {
	for _, a := range s {
		if a == u {
			return true
		}
	}
	return false
}
