use tokio::sync::broadcast;
use crate::models::state::user_manager::UserManager;

pub struct AppState<T>
where
    T: UserManager,
{
    user_state: T,
    transmitter: broadcast::Sender<String>
}