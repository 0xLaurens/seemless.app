#![allow(dead_code)]
use tracing::info;

#[cfg(test)]
mod tests;

mod app;
mod ws;
mod services;
mod models;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_target(false)
        .compact()
        .init();

    info!("Axum running 0.0.0.0:4000");

    axum::Server::bind(&"0.0.0.0:4000".parse().unwrap())
        .serve(app::create_app())
        .await
        .expect("Failed to launch app");
}
