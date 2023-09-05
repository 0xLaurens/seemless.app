use crate::models::user::{User, Username};


/*
* User management interface to easily change impl
* it's purpose is to manage the users connected and available to connect to.
*/
pub trait UserManager {
    fn add_user(&self, user: User) -> Option<User>;
    fn remove_user(&self, username: Username) -> Option<User>;
    fn update_user(&self, user: User, username: Username) -> Option<User>;
    fn get_users(&self) -> Vec<User>;
}