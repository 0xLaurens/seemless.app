use std::net::SocketAddr;
use std::sync::Arc;
use axum::extract::connect_info::IntoMakeServiceWithConnectInfo;
use axum::Router;
use axum::routing::get;
use crate::{services, ws};
use crate::models::state::app_state::AppState;
use crate::models::state::user_state_im::UserStateInMemory;

pub fn create_app() -> IntoMakeServiceWithConnectInfo<Router, SocketAddr> {
    let state = Arc::new(AppState::<UserStateInMemory>::new(UserStateInMemory::new()));
    create_routes(state)
        .into_make_service_with_connect_info::<SocketAddr>()
}

fn create_routes<S>(state: Arc<AppState<UserStateInMemory>>) -> Router<S> {
    Router::new()
        .merge(ws::websocket::create_routes(state.clone()))
        .route("/", get(|| async {":)"}))
        .layer(services::tracing::setup())
        .with_state(state)
}