use crate::models::user::{User, Username};


/*
* User management interface to easily change impl
* it's purpose is to manage the users connected and available to connect to.
*/
trait UserManager {
    fn add_user(&self, user: User) -> User;
    fn remove_user(&self, username: Username) -> User;
    fn update_user(&self, user: User, username: Username) -> User;
    fn get_users(&self) -> Vec<User>;
}