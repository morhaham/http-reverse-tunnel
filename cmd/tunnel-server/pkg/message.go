package pkg

type MessageType string

const (
	AuthenticationRequest  MessageType = "AUTHENTICATION_REQUEST"
	AuthenticationResponse MessageType = "AUTHENTICATION_RESPONSE"
	HttpProxy              MessageType = "HTTP_PROXY"
)

type Message struct {
	messageType MessageType
	content     string
}
