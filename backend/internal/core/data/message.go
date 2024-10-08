package data

import "github.com/gofiber/contrib/websocket"

type Message struct {
	Type        MessageType       `json:"type"`
	From        string            `json:"from,omitempty"`
	Target      string            `json:"target,omitempty"`
	User        *User             `json:"user,omitempty"`
	Users       []*User           `json:"users,omitempty"`
	SDP         string            `json:"sdp,omitempty"`
	Candidate   *RTCIceCandidate  `json:"candidate,omitempty"`
	Body        map[string]string `json:"body,omitempty"`
	RoomCode    RoomCode          `json:"roomCode,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	Conn        *websocket.Conn   `json:"-"`
}

type MessageType string

const (
	Offer           MessageType = "Offer"
	Answer          MessageType = "Answer"
	NewIceCandidate MessageType = "NewIceCandidate"

	PeerJoined  MessageType = "PeerJoined"
	PeerLeft    MessageType = "PeerLeft"
	PeerUpdated MessageType = "PeerUpdated"
	Peers       MessageType = "Peers"

	InvalidMessage MessageType = "InvalidMessage"

	RoomJoined      MessageType = "RoomJoined"
	RoomJoin        MessageType = "RoomJoin"
	RoomCreated     MessageType = "RoomCreated"
	RoomCodeInvalid MessageType = "RoomCodeInvalid"

	DisplayName          MessageType = "DisplayName"
	ChangeDisplayName    MessageType = "ChangeDisplayName"
	DuplicateDisplayName MessageType = "DuplicateDisplayName"
)
