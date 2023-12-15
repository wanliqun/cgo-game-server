package stress

import (
	"github.com/pkg/errors"
	"github.com/wanliqun/cgo-game-server/client"
)

type Rotbot struct {
	name   string
	client *client.Client
}

func (r *Rotbot) GetReady() error {
	err := r.client.Connect()
	if err != nil {
		return errors.WithMessage(err, "failed to connect client")
	}

	return nil
}

func (r *Rotbot) Close() {
	r.client.Close()
}
