package proto

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	invalidProtoMessage = errors.New("invalid protocol message")
)

func NewResponseMessage(msg proto.Message) *Message {
	resp := &Response{}
	switch v := msg.(type) {
	case *Status:
		resp.Body = &Response_Status{v}
	case *InfoResponse:
		resp.Body = &Response_Info{v}
	case *LoginResponse:
		resp.Body = &Response_Login{v}
	case *LogoutResponse:
		resp.Body = &Response_Logout{v}
	// TODO: extend for more message types support
	default:
		panic(invalidProtoMessage)
	}

	return &Message{Body: &Message_Response{resp}}
}
