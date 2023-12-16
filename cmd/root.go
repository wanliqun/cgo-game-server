package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/config"
)

var (
	configYaml string

	rootCmd = &cobra.Command{
		Use:   "cgo-game-server [--config | -c]",
		Short: "A demo high performance CGO game server.",
		Run:   run,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.InitKoanf(configYaml, cmd.Flags())
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configYaml,
		"config", "c", "config/config.yml",
		"YAML config file path to load",
	)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(simulatorCmd)
	rootCmd.AddCommand(loadRunnerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}
