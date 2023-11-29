package data

import "github.com/gofiber/contrib/websocket"

type Message struct {
	Type      MessageType       `json:"type"`
	From      string            `json:"from,omitempty"`
	Target    string            `json:"target,omitempty"`
	User      *User             `json:"user,omitempty"`
	Users     []*User           `json:"users,omitempty"`
	SDP       string            `json:"sdp,omitempty"`
	Candidate *RTCIceCandidate  `json:"candidate,omitempty"`
	Body      map[string]string `json:"body,omitempty"`
	Conn      *websocket.Conn   `json:"-"`
}

type MessageType string

var MessageTypes = struct {
	Offer             MessageType
	Answer            MessageType
	NewIceCandidate   MessageType
	PeerJoined        MessageType
	PeerLeft          MessageType
	PeerUpdated       MessageType
	Peers             MessageType
	Username          MessageType
	UsernamePrompt    MessageType
	InvalidMessage    MessageType
	DuplicateUsername MessageType
}{
	Offer:             "Offer",
	Answer:            "Answer",
	NewIceCandidate:   "NewIceCandidate",
	PeerJoined:        "PeerJoined",
	PeerLeft:          "PeerLeft",
	PeerUpdated:       "PeerUpdated",
	Peers:             "Peers",
	Username:          "Username",
	UsernamePrompt:    "UsernamePrompt",
	InvalidMessage:    "InvalidMessage",
	DuplicateUsername: "DuplicateUsername",
}
