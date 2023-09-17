package data

type Request struct {
	Type   RequestType `json:"type"`
	UserId UserID      `json:"userId"`
}

type RequestType string

const (
	Offer           RequestType = "Offer"
	Answer          RequestType = "Answer"
	NewIceCandidate RequestType = "NewIceCandidate"
	PeerJoined      RequestType = "PeerJoined"
	PeerLeft        RequestType = "PeerLeft"
	PeerUpdated     RequestType = "PeerUpdated"
)

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
