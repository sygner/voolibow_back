use chrono::Utc;
use serde::{Serialize, Deserialize};

#[derive(Serialize,Deserialize,Debug)]
pub struct User{
    pub phone_number:String,
    pub user_status:Status,
    pub user_id:i32,
    pub role:Role,
    pub created_at:chrono::DateTime<Utc>,
}

impl User{
    pub fn new(phone_number:&String)->Self{
        Self{
            user_id:0,
            phone_number:phone_number.to_owned(),
            user_status:Status::default(),
            role:Role::default(),
            created_at:chrono::Utc::now(),
        }
    }
}

#[derive(Debug,Serialize,Deserialize)]
pub enum Role{
    Owner,
    Admin,
    Moderator,
    User,
}

impl Default for Role{
    fn default() -> Self {
        Self::User
    }
}

impl ToString for Role{
    fn to_string(&self) -> String {
        match self{
            Role::Owner =>"Owner".to_owned(),
            Role::Admin => "Admin".to_owned(),
            Role::Moderator =>"Moderator".to_owned(),
            Role::User => "User".to_owned(),
        }
    }
}

impl From<String> for Role{
    fn from(value: String) -> Self {
        let value = value.as_str(); 
        match value{
            "Owner" => Role::Owner,
            "Admin" => Role::Admin,
            "Moderator" => Role::Moderator,
            "User" => Role::User,
            _ => Role::User,
        }
    }
}

#[derive(Debug,Serialize,Deserialize,PartialEq, Eq)]
pub enum Status{
    OnGoing(chrono::DateTime<Utc>),
    Suspended(chrono::DateTime<Utc>),
    Deleted(chrono::DateTime<Utc>),
    PermanentBan(chrono::DateTime<Utc>),
}

impl Default for Status{
    fn default() -> Self {
        Self::OnGoing(chrono::Utc::now())    
    }
}

impl ToString for Status{
    fn to_string(&self) -> String {
        match self{
            Status::OnGoing(e)=>format!("OnGoing {}",e),
            Status::Suspended(e)=>format!("Suspended {}",e),
            Status::Deleted(e)=>format!("Deleted {}",e),
            Status::PermanentBan(e)=>format!("PermanentBan {}",e)
        }
    }
}

impl From<String> for Status{
    fn from(value: String) -> Self {
        let (status,time) = value.split_once(' ').unwrap();

        match status{
            "OnGoing"=>Self::OnGoing(time.parse::<chrono::DateTime<Utc>>().unwrap()),
            "Suspended"=>Self::Suspended(time.parse::<chrono::DateTime<Utc>>().unwrap()),
            "Deleted"=>Self::Deleted(time.parse::<chrono::DateTime<Utc>>().unwrap()),
            _=>Self::OnGoing(time.parse::<chrono::DateTime<Utc>>().unwrap())
        }
    }
}

impl From<i32> for Status{
    fn from(value: i32) -> Self {
        match value{
            0=>Self::OnGoing(chrono::Utc::now()),
            1=>Self::Suspended(chrono::Utc::now()),
            2=>Self::Deleted(chrono::Utc::now()),
            3=>Self::PermanentBan(chrono::Utc::now()),
            _=>Self::OnGoing(chrono::Utc::now()),
        }
    }
}

impl Status{
    pub fn validate_status(&self)->bool{
        match &self{
            Status::OnGoing(_)=>true,
            _=>false
        }
    }
}
