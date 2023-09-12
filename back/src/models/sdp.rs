use serde::{Deserialize, Serialize};

/*
* SDP offers/answers contain the following fields
* https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling#exchanging_session_descriptions
*/
#[derive(Serialize, Deserialize, Debug)]
pub struct SessionDescriptionMessage {
    user: String,
    from: String,
    sdp: String,
}