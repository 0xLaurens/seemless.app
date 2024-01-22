export interface Connection {
    pc: RTCPeerConnection
    dc?: RTCDataChannel
    userId: string
}
