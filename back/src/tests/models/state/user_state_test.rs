use crate::models::state::user_state_im::UserStateInMemory;
use crate::models::user::User;
use crate::models::user_manager::UserManager;

#[test]
fn validate_new_user_state_is_empty() {
    let user_state = UserStateInMemory::new();
    assert_eq!(Vec::<User>::new(), user_state.get_users().unwrap());
}