import type {RoomCode} from "@/models/room";

export type Message = {
    type: string
    target?: string
    sdp?: string
    candidate?: RTCIceCandidateInit | null
    from?: string
    roomCode?: RoomCode
    displayName?: string
}
