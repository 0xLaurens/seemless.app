use tower_http::classify::{ServerErrorsAsFailures, SharedClassifier};
use tower_http::trace;
use tower_http::trace::TraceLayer;
use tracing::Level;

pub fn setup() -> TraceLayer<SharedClassifier<ServerErrorsAsFailures>> {
    TraceLayer::new_for_http()
        .make_span_with(trace::DefaultMakeSpan::new()
            .level(Level::INFO)
        )
        .on_response(trace::DefaultOnResponse::new()
            .level(Level::INFO)
        )
}