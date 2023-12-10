package proto

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	invalidProtoMessage = errors.New("invalid protocol message")
)

func NewResponseMessage(msg proto.Message) (*Message, error) {
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
	case *GenerateRandomNicknameResponse:
		resp.Body = &Response_GenerateRandomNickname{v}
	// TODO: extend for more message types support
	default:
		return nil, invalidProtoMessage
	}

	res := &Message{Body: &Message_Response{resp}}
	return res, nil
}

func NewRequestMessage(msg proto.Message) (*Message, error) {
	request := &Request{}
	switch v := msg.(type) {
	case *InfoRequest:
		request.Body = &Request_Info{v}
	case *LoginRequest:
		request.Body = &Request_Login{v}
	case *LogoutRequest:
		request.Body = &Request_Logout{v}
	case *GenerateRandomNicknameRequest:
		request.Body = &Request_GenerateRandomNickname{v}
	// TODO: extend for more message types support
	default:
		return nil, invalidProtoMessage
	}

	res := &Message{Body: &Message_Request{request}}
	return res, nil
}
