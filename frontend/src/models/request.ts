export enum RequestTypes {
    Offer = 'Offer',
    Answer = 'Answer',
    NewIceCandidate = 'NewIceCandidate',
    PeerJoined = 'PeerJoined',
    PeerLeft = 'PeerLeft',
    PeerUpdated = 'PeerUpdated',
    Peers = 'Peers',
    Username = 'Username',
    UsernamePrompt = 'UsernamePrompt',
    DuplicateUsername = 'DuplicateUsername',

    //New
    DisplayName = "DisplayName",


    PublicRoomLeft = "PublicRoomLeft",
    PublicRoomJoin = "PublicRoomJoin",
    PublicRoomPeers = "PublicRoomPeers",
    PublicRoomCreated = "PublicRoomCreated",
    PublicRoomIdInvalid = "PublicRoomIdInvalid",
    PublicRoomCreate = "PublicRoomCreate",

    JoinLocalRoom = "JoinLocalRoom",
    LeaveLocalRoom = "LeaveLocalRoom",
}
