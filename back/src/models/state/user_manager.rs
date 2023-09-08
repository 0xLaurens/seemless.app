use std::error::Error;
use crate::models::state::error::UserStateError;
use crate::models::user::{User, Username};


/*
* User management interface to easily change impl
* it's purpose is to manage the users connected and available to connect to.
*/
pub trait UserManager {
    fn add_user(&self, user: User) -> Result<Option<User>, UserStateError>;
    fn remove_user(&self, username: &Username) -> Result<Option<User>, UserStateError>;
    fn update_user(&self, user: User, username: Username) -> Result<Option<User>, UserStateError>;
    fn get_users(&self) -> Result<Vec<User>, Box<dyn Error>>;
}