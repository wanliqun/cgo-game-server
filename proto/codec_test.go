package proto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodecRequest(t *testing.T) {
	codec := NewCodec()

	msg, err := NewRequestMessage(&LoginRequest{
		Username: "kokko",
		Password: "helloworld",
	})
	assert.NoError(t, err, "failed to new request message")

	buf := bytes.NewBuffer(nil)

	err = codec.Encode(msg, buf)
	assert.NoError(t, err, "failed to encode message")

	msg2, err := codec.Decode(buf)
	assert.NoError(t, err, "failed to decode message")

	assert.Equal(
		t, msg.String(), msg2.String(), "unmarshalled message mismatched",
	)
}

func TestCodecResponse(t *testing.T) {
	codec := NewCodec()

	{
		msg, err := NewResponseMessage(&InfoResponse{
			ServerName:            "cgo-game-server",
			OnlinePlayers:         100,
			MaxPlayerCapacity:     1000,
			MaxConnectionCapacity: 10000,
		})
		assert.NoError(t, err, "failed to new response message")

		buf := bytes.NewBuffer(nil)

		err = codec.Encode(msg, buf)
		assert.NoError(t, err, "failed to encode message")

		msg2, err := codec.Decode(buf)
		assert.NoError(t, err, "failed to decode message")

		assert.Equal(
			t, msg.String(), msg2.String(), "unmarshalled message mismatched",
		)
	}

	{
		msg, err := NewResponseMessage(&GenerateRandomNicknameResponse{
			Nickname: "kokko",
		})
		assert.NoError(t, err, "failed to new response message")
		t.Log("msg", msg.String())

		buf := bytes.NewBuffer(nil)

		err = codec.Encode(msg, buf)
		assert.NoError(t, err, "failed to encode message")

		msg2, err := codec.Decode(buf)
		assert.NoError(t, err, "failed to decode message")

		t.Log("msg2", msg2.String())

		assert.Equal(
			t, msg.String(), msg2.String(), "unmarshalled message mismatched",
		)
	}
}
