use chrono::{Duration, Utc};

pub struct Token{
    pub user_id:i32,
    pub access_token:String,
    pub refresh_token:String,
    pub session_id:i16,
    pub role: String,
    pub agent:String,
    pub ip:String,
    pub status:Status,
    pub created_at:chrono::DateTime<Utc>,
    pub access_token_expire_at:chrono::DateTime<Utc>,
    pub refresh_token_expire_at:chrono::DateTime<Utc>,
}

impl Token{
    pub fn new(user_id:i32,agent:String,ip:String,role: String,expiry:i32) -> Self{
        let access_token = idgen::alpha_numeric(32); 
        let refresh_token = idgen::alpha_numeric(32);
        let session_id = idgen::numeric_code_i16(156,32767);
        Self{
            user_id,
            access_token,
            refresh_token,
            session_id,
            role,
            agent,
            ip,
            status:Status::Live,
            created_at:chrono::Utc::now(),
            access_token_expire_at:chrono::Utc::now() + Duration::try_seconds(expiry as i64).unwrap(),
            refresh_token_expire_at:chrono::Utc::now() + Duration::try_days(30).unwrap(),
        }
    }
    pub fn validate_status(&self)->bool{
        match self.status{
            Status::Banned=>false,
            Status::Live=>true
        }
    }

    pub fn validate_expiry(&self)->bool{
        self.access_token_expire_at > chrono::Utc::now()
    }
    
}

#[derive(PartialEq, Eq)]
pub enum Status{
    Banned,
    Live,
}

impl Default for Status{
    fn default() -> Self {
        Status::Live
    }
}

impl ToString for Status{
    fn to_string(&self) -> String {
        match self{
            Status::Banned=>format!("Banned"),
            Status::Live=>format!("Live")
        }   
    }
}

impl From<String> for Status{
    fn from(value: String) -> Self {
        let value = value.as_str();
        match value{
            "Banned"=>Self::Banned,
            "Live"=>Self::Live,
            _=>Self::Live
        }
    }
}

impl From<i32> for Status{
    fn from(value: i32) -> Self {
        match value{
            0=>Self::Banned,
            1=>Self::Live,
            _=>Self::Live,
        }
    }
}