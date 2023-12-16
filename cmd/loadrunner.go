package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/client"
	"github.com/wanliqun/cgo-game-server/stress"
)

type loadRunnerOption struct {
	*stress.ScheduleOption
	srvAddr string
	useUDP  bool
}

var (
	option = loadRunnerOption{
		ScheduleOption: &stress.ScheduleOption{},
	}

	loadRunnerCmd = &cobra.Command{
		Use:   "loadrunner",
		Short: "Load testing by producing high TPS",
		Run:   runLoadTest,
	}
)

func init() {
	loadRunnerCmd.Flags().IntVarP(
		&option.Workers,
		"workers", "w", 0, "Number of workers (go-routines)",
	)
	loadRunnerCmd.MarkFlagRequired("workers")

	loadRunnerCmd.Flags().IntVarP(
		&option.RobotsPerWorker,
		"robotsPerWorker", "r", 0, "Number of robots per worker",
	)
	loadRunnerCmd.MarkFlagRequired("robotsPerWorker")

	loadRunnerCmd.Flags().DurationVarP(
		&option.Interval,
		"iterationInterval", "i", 5*time.Second, "Iteration interval",
	)

	loadRunnerCmd.Flags().DurationVarP(
		&option.Timeout,
		"duration", "d", time.Minute, "Duration to produce TPS",
	)

	loadRunnerCmd.Flags().StringVarP(
		&option.srvAddr,
		"server-address", "s", "127.0.0.1:8765",
		"The address of server to be tested",
	)

	loadRunnerCmd.Flags().BoolVarP(
		&option.useUDP,
		"use-udp", "u", false,
		"Use UDP protocol to connect server",
	)
}

func runLoadTest(cmd *cobra.Command, args []string) {
	option.ClientFactory = func() *client.Client {
		if option.useUDP {
			return client.NewUDPClient(option.srvAddr)
		}

		return client.NewTCPClient(option.srvAddr)
	}

	scheduler := stress.NewScheduler(*option.ScheduleOption)
	scheduler.ProduceTPS(&stress.SimpleTask{})
}
