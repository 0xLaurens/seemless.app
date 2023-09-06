use std::collections::HashMap;
use std::error::Error;
use std::sync::{Arc, Mutex};
use crate::models::user::{User, Username};
use crate::models::user_manager::UserManager;

#[derive(Debug)]
pub struct UserStateInMemory {
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
    fn add_user(&self, user: User) -> Result<Option<User>, Box<dyn Error>> {
        let user = self.users
            .lock()
            .unwrap()
            .insert(user.get_username(), user);
        Ok(user)
    }

    fn remove_user(&self, username: &Username) -> Result<Option<User>, Box<dyn Error>> {
        let user = self.users
            .lock()
            .unwrap()
            .remove(username);
        Ok(user)
    }

    fn update_user(&self, user: User, username: Username) -> Result<Option<User>, Box<dyn Error>> {
        let _ = self.remove_user(&username);
        self.add_user(user)
    }

    fn get_users(&self) -> Result<Vec<User>, Box<dyn Error>> {
        let users = self.users
            .lock()
            .unwrap()
            .values()
            .map(|user| User::new(user.get_username(), user.get_user_agent()))
            .collect::<Vec<User>>();
        Ok(users)
    }
}