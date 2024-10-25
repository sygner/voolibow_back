use chrono::Utc;


pub struct LoginHistory{
    pub id:String,
    pub user_id:i32,
    pub user_role:String,
    pub section:String,
    pub ip:String,
    pub agent:String,
    pub logged_in_at:chrono::DateTime<Utc>
} 

impl LoginHistory{
    pub fn new(
        user_id:i32,
        user_role:String,
        section:String,
        ip:String,
        agent:String,
        logged_in_at:chrono::DateTime<Utc>
        )->Self{
        let id = idgen::alpha_numeric(30);
        Self{
            id,
            user_id,
            user_role,
            section,
            ip,
            agent,
            logged_in_at
        }
    }
}