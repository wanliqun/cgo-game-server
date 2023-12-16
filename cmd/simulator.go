package cmd

import (
	"log"
	"math/rand"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/client"
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/proto"
)

type simulatorOption struct {
	srvAddr  string
	userName string
	password string
}

var (
	rnd     *rand.Rand
	simOpts simulatorOption
	verbose bool

	simulatorCmd = &cobra.Command{
		Use:   "simulator",
		Short: "Simulates a game client to interact with server",
		Run:   runSimulator,
	}
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	simulatorCmd.Flags().StringVarP(
		&simOpts.srvAddr,
		"server-address", "s", "127.0.0.1:8765",
		"The server address this simulator connects to",
	)

	simulatorCmd.Flags().StringVarP(
		&simOpts.userName,
		"username", "u", "wanliqun",
		"The username used to login in the server",
	)

	simulatorCmd.Flags().StringVarP(
		&simOpts.password,
		"password", "p", "helloworld",
		"The password used to login in the server",
	)

	simulatorCmd.Flags().BoolVarP(
		&verbose,
		"verbose", "v", false,
		"Output debug log to the console",
	)
}

func runSimulator(*cobra.Command, []string) {
	if verbose {
		config.InitLogger(&config.LogConfig{
			Level: "debug", ForceColor: true,
		})
	}

	gc, err := chooseGameClient(simOpts.srvAddr)
	if err != nil {
		log.Fatalln("New game client error:", err)
	}
	defer gc.Close()

	gc.OnMessage(func(msg *proto.Message) {
		log.Println(">>> New message received from server:", msg.String())
	})

	if err := simulateGamePlay(gc); err != nil {
		log.Fatalln("Simulate game play error", err)
	}
}

func simulateGamePlay(gc *client.Client) (err error) {
	prompt := promptui.Select{
		Label: "Select Game Command",
		Items: []string{
			"INFO",
			"LOG IN",
			"LOG OUT",
			"GENERATE NICKNAME",
			"QUIT",
		},
	}

	for {
		idx, result, err := prompt.Run()
		if err != nil {
			return errors.WithMessage(err, "prompt error")
		}

		log.Println("Executing command:", result)

		switch idx {
		case 0: // info
			err = gc.Info()
		case 1: // login
			err = gc.Login(simOpts.userName, simOpts.password)
		case 2: // logout
			err = gc.Logout()
		case 3: // generate random nickname
			gender := common.Gender(rnd.Int() % 2)
			culture := common.Culture(rnd.Int() % 22)
			err = gc.GenerateRandomNickname(gender, culture)
		case 4: // quit
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func chooseGameClient(srvAddr string) (c *client.Client, err error) {
	prompt := promptui.Select{
		Label: "Select Client Type",
		Items: []string{"TCP", "UDP"},
	}

	idx, clientType, err := prompt.Run()
	if err != nil {
		return nil, errors.WithMessage(err, "prompt error")
	}

	log.Printf("You've chosen a %s client\n", clientType)
	if idx == 0 {
		c = client.NewTCPClient(srvAddr)
	} else {
		c = client.NewUDPClient(srvAddr)
	}

	if err := c.Connect(); err != nil {
		return nil, errors.WithMessage(err, "failed to connect client")
	}

	return c, nil
}
