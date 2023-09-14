package data

type Request struct {
	Type   RequestType `json:"type"`
	UserId UserID      `json:"userId"`
}

type RequestType int

const (
	Offer RequestType = iota
	Answer
	NewIceCandidate
	PeerJoined
	PeerLeft
	PeerUpdated
)

func (t RequestType) Parse() string {
	return []string{
		"Offer",
		"Answer",
		"NewIceCandidate",
		"PeerJoined",
		"PeerLeft",
		"PeerUpdated",
	}[t]
}
