package data

type Request struct {
	Type RequestType       `json:"type"`
	From string            `json:"from"`
	Body map[string]string `json:"body"`
}

type RequestType string

var RequestTypes = struct {
	Offer           RequestType
	Answer          RequestType
	NewIceCandidate RequestType
	PeerJoined      RequestType
	Peers           RequestType
	PeerLeft        RequestType
	PeerUpdated     RequestType
	Username        RequestType
}{
	Offer:           "Offer",
	Answer:          "Answer",
	NewIceCandidate: "NewIceCandidate",
	PeerJoined:      "PeerJoined",
	Peers:           "Peers",
	PeerLeft:        "PeerLeft",
	PeerUpdated:     "PeerUpdated",
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
