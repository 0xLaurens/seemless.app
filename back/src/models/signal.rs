/*
* different states of the application
* messages should be annotated with the type of signal
*/

use serde::{Deserialize, Serialize};
use crate::models::sdp::SessionDescriptionMessage;

#[derive(Serialize, Deserialize)]
#[serde(tag = "type")]
pub enum SignalType {
    Offer(SessionDescriptionMessage),
    Answer(SessionDescriptionMessage),
    NewIceCandidate,
    UserIdentity,
    PeerJoined,
    PeerLeft,
    Peers,
}