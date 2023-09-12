use std::sync::Arc;
use axum::extract::{State, WebSocketUpgrade};
use axum::extract::ws::{Message, WebSocket};
use axum::response::Response;
use axum::Router;
use axum::routing::get;
use futures::{SinkExt, StreamExt};
use futures::stream::{SplitSink, SplitStream};
use tokio::sync::broadcast::{Receiver, Sender};
use tokio::task::JoinHandle;
use crate::models::error::handler::HandlerError;
use crate::models::signal::SignalType;
use crate::models::state::app_state::AppState;
use crate::models::state::error::UserStateError;
use crate::models::state::user_manager::UserManager;
use crate::models::state::user_state_im::UserStateInMemory;
use crate::models::user::User;

pub fn create_routes<S>(state: Arc<AppState<UserStateInMemory>>) -> Router<S> {
    Router::new()
        .route("/ws/discover", get(discover_ws_incoming))
        .with_state(state)
}

async fn discover_ws_incoming(
    ws: WebSocketUpgrade,
    State(state): State<Arc<AppState<UserStateInMemory>>>,
) -> Response {
    ws.on_upgrade(move |socket| handle_discover_socket(socket, state))
}

async fn handle_discover_socket(
    socket: WebSocket,
    state: Arc<AppState<UserStateInMemory>>,
) {
    let (mut sender, mut receiver) = socket.split();

    let user = match receiver.next().await {
        Some(Ok(Message::Text(message))) => {
            setup_user(message, &state).await
        },
        _ => { Ok(None) },
    };

    dbg!(&user);

    let tx = state.get_transmitter();
    let rx = tx.subscribe();

    let _recv_task = handle_recv_task(rx, sender).await;
    let _conn_task = handle_connection_task(tx, receiver).await;
}

/*
* Task that keeps track of the websocket connection,
* send messages to broadcast channel
* and updates user state when the user disconnects
*/
async fn handle_connection_task(tx: Sender<String>, mut receiver: SplitStream<WebSocket>) -> JoinHandle<()> {
    tokio::spawn(async move {
        while let Some(Ok(Message::Text(msg))) = receiver.next().await {
            let _ = tx.send(msg);
        }

        let _ = tx.send(String::from("Left"));
    })
}

/*
* Task that handles incoming messages
* and calls other functions based on the message type
*/
async fn handle_recv_task(
    mut rx: Receiver<String>,
    mut sender: SplitSink<WebSocket, Message>,
) -> JoinHandle<()> {
    tokio::spawn(async move {
        while let Ok(msg) = rx.recv().await {
            dbg!(&msg);
            if sender.send(Message::Text(msg)).await.is_err() {
                break;
            }
        }
    })
}

/*
* Setup the user when they connect to the websocket for the first time;
*/
async fn setup_user(message: String, state: &AppState<UserStateInMemory>) -> Result<Option<User>, UserStateError> {
    match serde_json::from_str::<SignalType>(&message) {
        Ok(signal) => {
            match signal {
                SignalType::Join { username } => {
                    state.user_state.add_user(User::new(username, None))
                }
                _ => Ok(None)
            }
        }
        Err(_) => { Err(UserStateError::DeserializationError) }
    }
}

/*
* Match request type to the corresponding action
*/
async fn message_handler(message: String, state: &AppState<UserStateInMemory>) -> Result<(), HandlerError> {
    if let Ok(signal) = serde_json::from_str::<SignalType>(&message) {
        match signal {
            SignalType::Offer(_) => {
                Ok(())
            }
            SignalType::Answer(_) => {
                Ok(())
            }
            SignalType::NewIceCandidate => {
                Ok(())
            }
            SignalType::UserIdentity => {
                Ok(())
            }
            SignalType::PeerJoined => {
                Ok(())
            }
            SignalType::PeerLeft => {
                Ok(())
            }
            SignalType::Peers => {
                Ok(())
            }
            SignalType::Join { username } => {
                match state.user_state.add_user(User::new(username, None)) {
                    Ok(_) => {}
                    Err(_) => { return Err(HandlerError::UnexpectedError); }
                };
                Ok(())
            }
        }
    } else {
        Err(HandlerError::IncorrectMessageFormat)
    }
}