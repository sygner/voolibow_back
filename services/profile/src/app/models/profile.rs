pub struct Profile{
    pub user_id:i32,
    pub display_sid:String,
    pub display_username:String,
    pub username:String,
    pub avatar: String
}

impl Profile{
    pub fn new(user_id:i32,) ->Self{
        let display_sid = format!("{}:{}",idgen::alpha_numeric(32),idgen::numeric_code_i16(1231,4568));
        let display_username = idgen::alpha_numeric(20);
        let username = display_username.to_lowercase();
        Self{
            user_id,
            display_sid,
            display_username,
            username,
            avatar:String::new()
        }
    }
}