package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/client"
	"github.com/wanliqun/cgo-game-server/stress"
)

var (
	option stress.ScheduleOption

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
}

func runLoadTest(cmd *cobra.Command, args []string) {
	option.ClientFactory = func() *client.Client {
		//return client.NewTCPClient("127.0.0.1:8765")
		return client.NewUDPClient("192.168.2.152:8765")
	}

	scheduler := stress.NewScheduler(option)
	scheduler.ProduceTPS(&stress.SimpleTask{})
}
