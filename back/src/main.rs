mod app;
mod ws;

#[tokio::main]
async fn main() {
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app::create_app())
        .await
        .expect("Failed to launch app")
}
