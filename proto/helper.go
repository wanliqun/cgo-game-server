package proto

func NewStatusMessage(status *Status) *Message {
	return &Message{
		Body: &Message_Response{
			Response: &Response{
				Body: &Response_Status{
					Status: status,
				},
			},
		},
	}
}
