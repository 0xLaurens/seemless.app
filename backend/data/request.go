package data

type Request struct {
	Type   RequestType       `json:"type"`
	From   string            `json:"from"`
	Target string            `json:"target"`
	Body   map[string]string `json:"body"`
}

type RequestType string

var RequestTypes = struct {
	Offer           RequestType
	Answer          RequestType
	NewIceCandidate RequestType
	PeerJoined      RequestType
	PeerLeft        RequestType
	PeerUpdated     RequestType
	Peers           RequestType
	Username        RequestType
}{
	Offer:           "Offer",
	Answer:          "Answer",
	NewIceCandidate: "NewIceCandidate",
	PeerJoined:      "PeerJoined",
	PeerLeft:        "PeerLeft",
	PeerUpdated:     "PeerUpdated",
	Peers:           "Peers",
	Username:        "Username",
}

//func (t RequestType) Parse() string {
//	return []string{
//		"Offer",
//		"Answer",
//		"NewIceCandidate",
//		"PeerJoined",
//		"PeerLeft",
//		"PeerUpdated",
//	}[t]
//}
