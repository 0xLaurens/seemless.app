export interface Connection {
  pc: RTCPeerConnection
  dc?: RTCDataChannel
  username: string
}
