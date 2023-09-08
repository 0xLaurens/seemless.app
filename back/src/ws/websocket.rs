use axum::extract::WebSocketUpgrade;
use axum::extract::ws::WebSocket;
use axum::response::Response;
use axum::Router;
use axum::routing::get;
use futures::{StreamExt};

pub fn create_routes() -> Router {
    Router::new()
        .route("/ws/discover", get(discover_ws_incoming))
}

async fn discover_ws_incoming(
    ws: WebSocketUpgrade,
) -> Response {
    ws.on_upgrade(handle_discover_socket)
}

async fn handle_discover_socket(
    socket: WebSocket
) {
    let (mut _sender, _receiver) = socket.split();
}

/*
* Task that keeps track of the websocket connection,
* broadcasts requests sent the channel
* and updates user state when the user disconnects
*/
async fn handle_connection_task()  {
    todo!()
}

/*
* Task that handles incoming messages
* and calls other functions based on the message type
*/
async fn handle_recv_task() {
   todo!()
}