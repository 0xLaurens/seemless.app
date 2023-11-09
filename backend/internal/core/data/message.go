package data

type Message struct {
	Type      MessageType       `json:"type"`
	From      string            `json:"from"`
	Target    string            `json:"target"`
	User      *User             `json:"user"`
	Users     []*User           `json:"users"`
	SDP       string            `json:"sdp"`
	Candidate *RTCIceCandidate  `json:"candidate"`
	Body      map[string]string `json:"body"`
}

type MessageType string

var MessageTypes = struct {
	Offer           MessageType
	Answer          MessageType
	NewIceCandidate MessageType
	PeerJoined      MessageType
	PeerLeft        MessageType
	PeerUpdated     MessageType
	Peers           MessageType
	Username        MessageType
	UsernamePrompt  MessageType
	InvalidMessage  MessageType
}{
	Offer:           "Offer",
	Answer:          "Answer",
	NewIceCandidate: "NewIceCandidate",
	PeerJoined:      "PeerJoined",
	PeerLeft:        "PeerLeft",
	PeerUpdated:     "PeerUpdated",
	Peers:           "Peers",
	Username:        "Username",
	UsernamePrompt:  "UsernamePrompt",
	InvalidMessage:  "InvalidMessage",
}
