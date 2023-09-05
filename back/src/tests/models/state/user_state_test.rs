use crate::models::state::user_state_im::UserStateInMemory;
use crate::models::user::User;
use crate::models::user_manager::UserManager;

#[test]
fn validate_new_user_state_is_empty() {
    let user_state = UserStateInMemory::new();
    assert_eq!(Vec::<User>::new(), user_state.get_users().unwrap());
}

#[test]
fn validate_get_one_user() {
    let user_state = UserStateInMemory::new();
    let user = User::new(String::from("JohnyTest"), None);
    let _ = user_state.add_user(user.clone());
    assert_eq!(vec![user], user_state.get_users().unwrap());
}

#[test]
fn validate_get_three_users() {
    let user_state = UserStateInMemory::new();
    let users = vec![
        User::new(String::from("JohnyTest"), None),
        User::new(String::from("Gary"), None),
        User::new(String::from("Fritz"), None)
    ];

    for user in users.clone() {
        let _ = user_state.add_user(user);
    }

    assert_eq!(user_state.get_users().unwrap().len(), 3);
    assert!(user_state.get_users().unwrap().contains(&users[0]));
    assert!(user_state.get_users().unwrap().contains(&users[1]));
    assert!(user_state.get_users().unwrap().contains(&users[2]));
}