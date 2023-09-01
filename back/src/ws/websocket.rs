use axum::extract::WebSocketUpgrade;
use axum::extract::ws::WebSocket;
use axum::response::Response;
use axum::Router;
use axum::routing::get;
use futures::{SinkExt, StreamExt};

pub fn create_routes() -> Router {
    Router::new()
        .route("/ws/discover", get(discover_ws_incoming))
}

async fn discover_ws_incoming(
    ws: WebSocketUpgrade,
) -> Response {
    ws.on_upgrade(move |socket| handle_discover_socket(socket))
}

async fn handle_discover_socket(
    socket: WebSocket
) {
    let (mut sender, mut _receiver) = socket.split();
    if sender.send("Welcome".clone().into()).await.is_ok() {
        println!("Sent: Welcome message")
    }

}