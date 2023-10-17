export type Message = {
  type: string
  target?: string
  sdp?: string
  candidate?: RTCIceCandidateInit | null
  from?: string
}
