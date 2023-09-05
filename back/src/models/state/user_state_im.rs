use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use crate::models::user::{User, Username};

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