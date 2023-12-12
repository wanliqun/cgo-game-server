package cmd

import (
	"log"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanliqun/cgo-game-server/client"
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/proto"
)

const (
	srvAddr  = "127.0.0.1:8765"
	username = "kokko"
	password = "helloworld"
	sex      = common.Male
	culture  = common.CHINESE
)

var (
	simulatorCmd = &cobra.Command{
		Use:   "simulator",
		Short: "Simulates a game client to interact with server",
		Run:   runSimulator,
	}
)

func runSimulator(*cobra.Command, []string) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})

	gc, err := chooseGameClient(srvAddr)
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
			err = gc.Login(username, password)
		case 2: // logout
			err = gc.Logout()
		case 3: // generate random nickname
			err = gc.GenerateRandomNickname(sex, culture)
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
