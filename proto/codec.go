package proto

import (
	"encoding/binary"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// Codec serializes and deserializes data between protocol message and underlying network package data.
type Codec struct {
	proto.MarshalOptions
	proto.UnmarshalOptions
}

func NewCodec() *Codec {
	return &Codec{
		MarshalOptions: proto.MarshalOptions{
			Deterministic: true, // use deterministic ordering for map fields
		},
		UnmarshalOptions: proto.UnmarshalOptions{
			DiscardUnknown: true, // discard unknown fields
		},
	}
}

func (c *Codec) Encode(msg *Message, w io.Writer) error {
	data, err := c.Marshal(msg)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal message")
	}

	// Write message length as the first 4 bytes in big endian.
	if err := binary.Write(w, binary.BigEndian, int32(len(data))); err != nil {
		return errors.WithMessage(err, "failed to write message length")
	}

	// Write message data.
	if _, err := w.Write(data); err != nil {
		return errors.WithMessage(err, "failed to write msg data")
	}

	logrus.WithFields(logrus.Fields{
		"dataLen": len(data),
		"dataHex": hex.EncodeToString(data),
		"msg":     msg.String(),
	}).Debug("Codec encodes message")

	return nil
}

func (c *Codec) Decode(r io.Reader) (*Message, error) {
	// Read message length.
	len := int32(0)
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, errors.WithMessage(err, "failed to read message length")
	}

	// Read message data.
	data := make([]byte, len)
	if _, err := r.Read(data); err != nil {
		return nil, errors.WithMessage(err, "failed to read message data")
	}

	msg := new(Message)
	if err := c.Unmarshal(data, msg); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal msg data")
	}

	logrus.WithFields(logrus.Fields{
		"dataLen": len,
		"dataHex": hex.EncodeToString(data),
		"msg":     msg.String(),
	}).Debug("Codec decodes message")

	return msg, nil
}
