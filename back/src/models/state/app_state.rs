use tokio::sync::broadcast;
use crate::models::state::user_manager::UserManager;

pub struct AppState<T>
where
    T: UserManager + Clone,
{
    pub user_state: T,
    transmitter: broadcast::Sender<String>
}

impl<T: UserManager + Clone> AppState<T> {
    pub fn new(user_state: T) -> Self {
       Self {
           user_state,
           transmitter: broadcast::channel(420).0,
       }
    }

    pub fn get_transmitter(&self) -> broadcast::Sender<String> {
        self.transmitter.clone()
    }
}