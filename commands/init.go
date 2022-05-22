package commands

import (
	"Noctobot/commands/admin"
	"Noctobot/commands/config"
	"Noctobot/commands/game"
	"Noctobot/commands/hypixel"
	"Noctobot/commands/info"
	"Noctobot/commands/ksoft"
	"Noctobot/commands/mod"
	"Noctobot/commands/util"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
)

func Init(bot *gocto.Bot) {
	// Owner
	bot.AddCommand(gocto.NewCommand("js", "Owner", admin.OwnerEval).SetDescription("Evaluate JavaScript.").SetOwnerOnly(true).SetEditable(true).SetUsage("<code:string...>").AddAliases("ev"))
	bot.AddCommand(gocto.NewCommand("py", "Owner", admin.OwnerStarlark).SetDescription("Evaluate JavaScript.").SetOwnerOnly(true).SetUsage("<code:string...>").AddAliases("ev2"))
	bot.AddCommand(gocto.NewCommand("go", "Owner", admin.OwnerAnko).SetDescription("Evaluates arbitrary Anko").SetOwnerOnly(true).SetUsage("<code:string...>"))
	bot.AddCommand(gocto.NewCommand("su", "Owner", admin.OwnerSu).SetDescription("Run a command as another user.").SetOwnerOnly(true).SetUsage("<@user> <command:string> [args:string...]"))
	bot.AddCommand(gocto.NewCommand("sudo", "Owner", admin.OwnerSudo).SetDescription("Run a command without cooldown.").SetOwnerOnly(true).SetUsage(" <command:string> [args:string...]"))
	bot.AddCommand(gocto.NewCommand("exec", "Owner", admin.AdminExec).SetDescription("Executes shell commands.").SetOwnerOnly(true).SetUsage("<command:string...>"))
	bot.AddCommand(gocto.NewCommand("debug", "Owner", admin.Debug).SetOwnerOnly(true).SetUsage("<cmd:string> [args:string...]").NoOverride(false))
	bot.AddCommand(gocto.NewCommand("sql", "Owner", admin.OwnerSQL).SetDescription("Execute some SQL queries.").SetOwnerOnly(true).SetUsage("<query:string...>"))
	bot.AddCommand(gocto.NewCommand("lg", "Owner", admin.Listguilds).SetUsage("[sorting:string]").SetOwnerOnly(true))
	bot.AddCommand(gocto.NewCommand("shutdown", "Owner", admin.OwnerShutdown).SetDescription("Shuts down the bot.").SetOwnerOnly(true).AddAliases("reboot"))

	bot.AddCommand(gocto.NewCommand("fetch", "Owner", admin.OwnerFetch).SetOwnerOnly(true).SetUsage("<type:string> <id:int>"))

	// Admin
	bot.AddCommand(gocto.NewCommand("prefix", "Config", config.Prefix).SetDescription("Get the current prefix for this guild").SetGuildOnly(true).SetUsage("[prefix:string]").SetGuildOnly(true).AddAliases("setprefix", "changeprefix"))
	bot.AddCommand(gocto.NewCommand("setprefix", "Config", config.SetPrefix).SetDescription("Sets the prefix for this server.\nSurround with `\"` to save spaces").SetGuildOnly(true).SetPermission(discordgo.PermissionManageServer).SetUsage("[prefix:string]").SetGuildOnly(true).AddAliases("setprefix", "changeprefix"))
	bot.AddCommand(gocto.NewCommand("config", "Config", config.Config).SetDescription("Change a config key's value").SetGuildOnly(true).SetPermission(discordgo.PermissionManageServer).SetUsage("<key:string> <value:string...>").SetGuildOnly(true).AddAliases("setprefix", "changeprefix"))

	// Moderators
	bot.AddCommand(gocto.NewCommand("purge", "Moderation", mod.Purge).SetUsage("<count:int> [channel:channel] [@author]").SetGuildOnly(true).SetDescription("Delete messages").SetPermission(discordgo.PermissionManageMessages).AddAliases("prune").Delete())
	bot.AddCommand(gocto.NewCommand("warn", "Moderation", mod.Warn).SetUsage("<@user> [reason:string]").SetDescription("Warn someone").SetGuildOnly(true).SetPermission(discordgo.PermissionManageMessages).Delete())
	bot.AddCommand(gocto.NewCommand("warns", "Moderation", mod.Warns).SetUsage("<@user>").SetDescription("See someone's warns").SetGuildOnly(true).SetPermission(discordgo.PermissionManageMessages).AddAliases("warnings"))
	bot.AddCommand(gocto.NewCommand("kick", "Moderation", mod.Kick).SetUsage("<@user> <reason:string>").SetDescription("Gone 410").SetGuildOnly(true).SetPermission(discordgo.PermissionKickMembers).Delete())
	bot.AddCommand(gocto.NewCommand("ban", "Moderation", mod.Ban).SetUsage("<@user> <reason:string>").SetDescription("Banish someone").SetGuildOnly(true).SetPermission(discordgo.PermissionBanMembers).Delete())
	bot.AddCommand(gocto.NewCommand("clearwarns", "Moderation", mod.Clearwarns).SetUsage("<@user>").SetDescription("Clear someone's warns").SetGuildOnly(true).SetPermission(discordgo.PermissionManageMessages).AddAliases("clearwarn", "clearnwarning", "clearnwarnings"))
	bot.AddCommand(gocto.NewCommand("deletecase", "Moderation", mod.DeleteCase).SetDescription("Remove a case from someone's profile").SetUsage("<id:int>").SetGuildOnly(true).SetPermission(discordgo.PermissionManageMessages).AddAliases("deletewarn", "deletewarning", "removecase", "removewarn", "removewarning"))
	//bot.AddCommand(gocto.NewCommand("unban", "Moderation", mod.Unban).SetUsage("<@user>").SetPermission(discordgo.PermissionBanMembers).Delete())

	// Info
	bot.AddCommand(gocto.NewCommand("userinfo", "Information", info.User).SetUsage("[@user]").AddAliases("ui"))
	bot.AddCommand(gocto.NewCommand("serverinfo", "Information", info.Server).SetGuildOnly(true).SetUsage("[id:int]").AddAliases("si", "gi"))

	// Hypixel
	bot.AddCommand(gocto.NewCommand("hypixel", "Hypixel", hypixel.Hypixel).SetDescription("View someone's stats on hypixel.net").SetUsage("<target:string>"))

	bot.AddCommand(gocto.NewCommand("tictactoe", "Games", game.TicTacToe).SetGuildOnly(true).AddAliases("ttt"))

	// Lyrics
	bot.AddCommand(gocto.NewCommand("lyrics", "KSoft", ksoft.Lyrics).SetDescription("View your favourite's music").SetUsage("[song:string]"))
	bot.AddCommand(gocto.NewCommand("baninfo", "KSoft", ksoft.BanInfo).SetUsage("<target:string>").AddAliases("bi"))
	bot.AddCommand(gocto.NewCommand("gis", "KSoft", ksoft.GIS).SetUsage("<location:string...>"))
	bot.AddCommand(gocto.NewCommand("weather", "KSoft", ksoft.Weather).SetUsage("<location:string...>").AddAliases("w"))
	bot.AddCommand(gocto.NewCommand("meme", "KSoft", ksoft.Meme).SetDescription("Get a random meme"))

	// Utility
	bot.AddCommand(gocto.NewCommand("colour", "Utility", util.Color).SetDescription("View a color by it's code").SetUsage("<hexcolour:string>").AddAliases("color"))
	bot.AddCommand(gocto.NewCommand("attach", "Utility", util.ReadAttachment).SetDescription("View a message's attachment").SetUsage("<msgid:string>"))
	bot.AddCommand(gocto.NewCommand("avatar", "Utility", util.Avatar).SetDescription("View someone's avatar.").SetUsage("[@user]").AddAliases("av", "pp"))
	bot.AddCommand(gocto.NewCommand("qrcode", "Utility", util.QRCode).SetDescription("QRCode utility").SetUsage("<type:string> [content:string...]").AddAliases("qr").SetCooldown(5))

	bot.AddCommand(gocto.NewCommand("suggest", "General", util.Suggest).SetDescription("Suggest a feature").SetUsage("<content:string...>").SetCooldown(30))
	bot.AddCommand(gocto.NewCommand("support", "General", util.Support).SetDescription("Discord server to get support").SetCooldown(60))
}
