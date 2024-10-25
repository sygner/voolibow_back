use super::types::error::Error;

pub fn check_username(display_username:String)->Result<(),Error>{
    let rgx = regex::Regex::new(r#"^[a-zA-Z]+$"#).map_err(|_| Error::InternalError("failed to validate".to_owned()))?;
    if !rgx.is_match(&display_username){
        return Err(Error::ServiceError("the format of username is incorrect #111".to_owned()))
    };

    Ok(())
}