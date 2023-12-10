package server

import (
	"context"

	"github.com/wanliqun/cgo-game-server/proto"
)

type HandlerFunc func(context.Context, *Message) *Message

type Message struct {
	*proto.Message
	Error error
}

func NewMessage(msg *proto.Message) *Message {
	return &Message{Message: msg}
}

func NewMessageWithError(err error) *Message {
	if err == nil {
		err = NilError
	}
	return &Message{Error: err}
}

// ProtoMessage adapts into a protobuf response message.
func (m *Message) ProtoMessage() *proto.Message {
	if m.Error == nil {
		return m.Message
	}

	status := StatusInternalServerError
	if v, ok := m.Error.(Error); ok {
		status = v.Status()
	}

	return proto.NewResponseMessage(&proto.Status{
		Code:    status,
		Message: m.Error.Error(),
	})
}
