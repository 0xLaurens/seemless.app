use std::error::Error;
use crate::models::state::user_state_im::UserStateInMemory;
use crate::models::user::{User, Username};
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

#[test]
fn remove_non_existing_user() {
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
    let _ = user_state.remove_user(&Username::from("Johny"));
    assert_eq!(user_state.get_users().unwrap().len(), 3);
}

#[test]
fn remove_one_user() {
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
    let _ = user_state.remove_user(&users[0].get_username());
    assert_eq!(user_state.get_users().unwrap().len(), 2);
}

#[test]
fn update_user_list_contains_new_username() -> Result<(), Box<dyn Error>> {
    let user_state = UserStateInMemory::new();
    let users = vec![
        User::new(String::from("JohnyTest"), None),
        User::new(String::from("Gary"), None),
        User::new(String::from("Fritz"), None)
    ];

    for user in users.clone() {
        let _ = user_state.add_user(user);
    }
    assert_eq!(user_state.get_users()?.len(), 3);
    assert!(user_state.get_users()?.contains(&users[0]));


    assert_eq!(user_state.get_users()?.contains(&users[0]), true);

    let updated_user = User::new(String::from("John"), None);
    let _ = user_state.update_user(updated_user.clone(), users[0].get_username());
    assert_ne!(user_state.get_users()?.contains(&users[0]), true);
    assert_eq!(user_state.get_users()?.contains(&updated_user), true);
    Ok(())
}