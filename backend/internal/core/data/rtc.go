package data

type RTCIceCandidate struct {
	Candidate        string `json:"candidate"`
	SdpMid           string `json:"sdpMid"`
	SdpMLineIndex    int    `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment"`
}
