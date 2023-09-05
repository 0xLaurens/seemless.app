#![allow(dead_code)]
use std::net::SocketAddr;
use axum::extract::connect_info::IntoMakeServiceWithConnectInfo;
use axum::Router;
use axum::routing::get;
use crate::{services, ws};

pub fn create_app() -> IntoMakeServiceWithConnectInfo<Router, SocketAddr>
{
    Router::new()
        .merge(ws::websocket::create_routes())
        .route("/", get(|| async {":)"}))
        .layer(services::tracing::setup())
        .into_make_service_with_connect_info::<SocketAddr>()
}