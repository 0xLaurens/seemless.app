/*
* different states of the application
* messages should be annotated with the type of signal
*/

use serde::{Deserialize, Serialize};
use crate::models::sdp::SessionDescriptionMessage;
use crate::models::user::Username;

#[derive(Serialize, Deserialize, Debug)]
#[serde(tag = "type")]
pub enum SignalType {
    Offer(SessionDescriptionMessage),
    Answer(SessionDescriptionMessage),
    NewIceCandidate,
    UserIdentity,
    PeerJoined,
    PeerLeft,
    Peers,
    Join {
        username: Username
    },
}