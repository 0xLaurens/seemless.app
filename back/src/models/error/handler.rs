use std::fmt;
use std::fmt::{Display, Formatter};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub enum HandlerError {
    IncorrectMessageFormat,
    UnexpectedError,
}

impl Display for HandlerError {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match *self {
            HandlerError::IncorrectMessageFormat => write!(f, "Incorrect Request Format"),
            _ => write!(f, "Unexpected Error")
        }
    }
}