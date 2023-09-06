use std::fmt::Display;
use std::error::Error;
use std::fmt;
use std::fmt::Formatter;
use std::sync::PoisonError;

#[derive(Debug, PartialEq)]
pub enum UserStateError {
    UsernameNotUnique,
    UserNotFound,
    PoisonError,
}

impl Display for UserStateError {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match *self {
            UserStateError::UsernameNotUnique =>
                write!(f, "The username is not unique, please select another username"),
            UserStateError::UserNotFound =>
                write!(f, "Could not find the user matching the specified username"),
            _ => write!(f, "Unexpected error"),
        }
    }
}

impl Error for UserStateError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match *self {
            UserStateError::UsernameNotUnique => None,
            UserStateError::UserNotFound => None,
            _ => None
        }
    }
}

impl<T> From<PoisonError<T>> for UserStateError {
    fn from(_err: PoisonError<T>) -> Self {
        Self::PoisonError
    }
}