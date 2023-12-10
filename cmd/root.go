package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/game"
)

var (
	configYaml string

	rootCmd = &cobra.Command{
		Use:   "cgo-game-server [--config | -c]",
		Short: "A demo high performance CGO game server.",
		Run:   run,
	}
)

func init() {
	rootCmd.Flags().StringVarP(
		&configYaml,
		"config", "c", "config/config.yml",
		"YAML config file path to load",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	app, err := game.NewApplication(configYaml)
	if err != nil {
		panic(err)
	}

	app.Run()
}
