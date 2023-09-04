/*
* different states of the application
* messages should be annotated with the type of signal
*/

pub enum SignalType {
    Offer,
    Answer,
    NewIceCandidate,
    UserIdentity,
    PeerJoined,
    PeerLeft,
    Peers,
}