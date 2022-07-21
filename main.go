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

	flag.Parse()
	ctx := context.Background()

	os.Exit(int(subcommands.Execute(ctx)))
}
