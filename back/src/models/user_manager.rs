use crate::models::user::{User, Username};

trait UserManager {
    fn add_user(&self, user: User) -> User;
    fn remove_user(&self, username: Username) -> User;
    fn update_user(&self, user: User, username: Username) -> User;
    fn get_users(&self) -> Vec<User>;
}