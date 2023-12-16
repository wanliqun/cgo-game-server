package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/game"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start game server",
		Run:   runServer,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.InitKoanfFromPflag(cmd.Flags())
		},
	}
)

func init() {
	serverCmd.Flags().Bool(
		"cgo.enabled", false,
		"Whether to enable CGO",
	)

	serverCmd.Flags().String(
		"cgo.resourceDir", "./resources",
		"Resource path for CGO monicker generator",
	)
}

func runServer(cmd *cobra.Command, args []string) {
	app, err := game.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Run()
}
