pub type Username = String;

#[derive(Debug, PartialEq)]
pub struct User {
 username: Username,
 user_agent: Option<String>,
}

impl User {
    pub fn new(username: Username, user_agent: Option<String>) -> Self {
        Self {
           username,
           user_agent,
       }
    }

    pub fn get_username(&self) -> Username {
        self.username.clone()
    }

    pub fn get_user_agent(&self) -> Option<String> {
        self.user_agent.clone()
    }
}