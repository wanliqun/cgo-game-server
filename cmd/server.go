package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/game"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start game server",
		Run:   runServer,
	}
)

func runServer(*cobra.Command, []string) {
	app, err := game.NewApplication(configYaml)
	if err != nil {
		panic(err)
	}

	app.Run()
}
