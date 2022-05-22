package admin

import (
	"Noctobot/database"
	"github.com/Noctember/gocto"
	"os"
)

func OwnerShutdown(ctx *gocto.CommandContext) {
	ctx.Reply("Shutting down...")

	ctx.Session.Close()
	database.Close()

	os.Exit(0)
}
