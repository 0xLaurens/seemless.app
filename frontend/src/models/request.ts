export enum RequestTypes {
    Offer = 'Offer',
    Answer = 'Answer',
    NewIceCandidate = 'NewIceCandidate',

    PeerJoined = 'PeerJoined',
    PeerLeft = 'PeerLeft',
    PeerUpdated = 'PeerUpdated',
    Peers = 'Peers',

    RoomJoined = "RoomJoined",
    RoomJoin = "RoomJoin",
    RoomCreated = "RoomCreated",
    RoomCodeInvalid = "RoomCodeInvalid",

    DisplayName = "DisplayName",
    ChangeDisplayName = "ChangeDisplayName",
    DuplicateUsername = "DuplicateUsername",
}
