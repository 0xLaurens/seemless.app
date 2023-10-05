export interface SessionDescriptionMessage extends RTCSessionDescriptionInit {
  from: string
  target: string
  sdp?: string | undefined
  type: RTCSdpType
}
