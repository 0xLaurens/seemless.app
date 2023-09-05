use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tracing::Instrument;
use crate::models::user::{User, Username};
use crate::models::user_manager::UserManager;

struct UserStateInMemory {
    users: Arc<Mutex<HashMap<Username, User>>>
}

impl UserStateInMemory {
    pub fn new() -> Self {
        Self {
            users: Arc::new(Mutex::new(HashMap::new()))
        }
    }
}

impl UserManager for UserStateInMemory {
    fn add_user(&self, user: User) -> Option<User> {
        self.users.lock()
            .insert(user.get_username(), user.clone())
    }

    fn remove_user(&self, username: Username) -> User {
        self.users
            .lock()
            .remove(username)?
    }

    fn update_user(&self, user: User, username: Username) -> Option<User> {
       let _ = self.users
            .lock()?
            .remove(&username);

        self.users
            .lock()?
            .insert(username, user)
    }

    fn get_users(&self) -> Vec<User> {
        self.users
            .clone()
            .values()
            .collect()
    }
}