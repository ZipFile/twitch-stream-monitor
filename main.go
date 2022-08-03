package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	_ "github.com/joho/godotenv/autoload"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
	tsm_cli "github.com/ZipFile/twitch-stream-monitor/internal/cli"
)

func main() {
	app := tsm_app.New()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(tsm_cli.NewMonitor(app), "")
	subcommands.Register(tsm_cli.NewSubscribe(app), "")
	subcommands.Register(tsm_cli.NewUnsubscribe(app), "")
	subcommands.Register(tsm_cli.NewGetAppAccessToken(app), "")
	subcommands.Register(tsm_cli.NewResolveUsername(app), "")
	subcommands.Register(tsm_cli.NewList(app), "")

	flag.Parse()
	ctx := context.Background()

	os.Exit(int(subcommands.Execute(ctx)))
}
