use tokio::sync::broadcast;

pub struct AppState<T>
where
    T: UserManager,
{
    user_state: T,
    transmitter: broadcast::Sender<String>
}