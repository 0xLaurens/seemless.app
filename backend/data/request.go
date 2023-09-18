package data

type Request struct {
	Type   RequestType `json:"type"`
	UserId UserID      `json:"userId"`
}

type RequestType string

var RequestTypes = struct {
	Offer           RequestType
	Answer          RequestType
	NewIceCandidate RequestType
	PeerJoined      RequestType
	PeerLeft        RequestType
	PeerUpdated     RequestType
}{
	Offer:           "Offer",
	Answer:          "Answer",
	NewIceCandidate: "NewIceCandidate",
	PeerJoined:      "PeerJoined",
	PeerLeft:        "PeerLeft",
	PeerUpdated:     "PeerUpdated",
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
