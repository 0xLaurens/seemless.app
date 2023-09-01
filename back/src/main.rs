mod app;
mod ws;
mod services;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_target(false)
        .compact()
        .init();

    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app::create_app())
        .await
        .expect("Failed to launch app")
}
