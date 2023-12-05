package proto

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var ()

// Codec serializes and deserializes data between protocol message and underlying network package data.
type Codec struct {
	proto.MarshalOptions
	proto.UnmarshalOptions
}

func NewCodec() *Codec {
	return &Codec{
		MarshalOptions: proto.MarshalOptions{
			Deterministic: true, // use deterministic ordering for map fields
			AllowPartial:  true, // allow marshaling messages with missing required fields
			UseCachedSize: true, // use the cached size of the message if available
		},
		UnmarshalOptions: proto.UnmarshalOptions{
			AllowPartial:   true,  // allow unmarshaling messages with missing required fields
			DiscardUnknown: true,  // discard unknown fields
			Merge:          false, // do not merge with existing message
		},
	}
}

func (c *Codec) Encode(msg proto.Message, w io.Writer) error {
	data, err := c.Marshal(msg)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal message")
	}

	// Write message length as the first 4 bytes in big endian.
	if err := binary.Write(w, binary.BigEndian, c.Size(msg)); err != nil {
		return errors.WithMessage(err, "failed to write message length")
	}

	// Write message data.
	if _, err := w.Write(data); err != nil {
		return errors.WithMessage(err, "failed to write msg data")
	}

	return nil
}

func (c *Codec) Decode(r io.Reader) (*Message, error) {
	// Read message length.
	len := 0
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, errors.WithMessage(err, "failed to read message length")
	}

	// Read message data.
	data := make([]byte, 0, len)
	if _, err := r.Read(data); err != nil {
		return nil, errors.WithMessage(err, "failed to read message data")
	}

	msg := new(Message)
	if err := c.Unmarshal(data, msg); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal msg data")
	}

	return msg, nil
}
