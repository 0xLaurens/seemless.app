use std::fmt::Display;
use std::error::Error;
use std::fmt;
use std::fmt::Formatter;
use std::sync::PoisonError;

#[derive(Debug)]
pub enum UserStateError {
    UsernameNotUnique,
    UserNotFound
}

impl Display for UserStateError {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match *self {
            UserStateError::UsernameNotUnique =>
                write!(f, "The username is not unique, please select another username"),
            UserStateError::UserNotFound =>
                write!(f, "Could not find the user matching the specified username"),
        }
    }
}

impl Error for UserStateError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match *self {
            UserStateError::UsernameNotUnique => None,
            UserStateError::UserNotFound => None
        }
    }
}

impl<T> From<PoisonError<T>> for UserStateError {
    fn from(err: PoisonError<T>) -> Self {
        Self::from(err)
    }
}