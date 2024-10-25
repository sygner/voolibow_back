use reqwest::StatusCode;
use serde::{Serialize,Deserialize};

use super::types::error::Error;


const API_KEY:&str = "_";

#[derive(Serialize,Deserialize)]
pub struct SMS{
    pub message:String,
    pub line_number:String,
    pub phone_numbers:Vec<String>,
}

impl SMS{
    pub fn new_message(message:String,line_number:String,phone_numbers:Vec<String>)->Self{
        Self{
            message,
            line_number,
            phone_numbers,
        }
    }
    pub fn new_verification_message(code:i16,line_number:String,phone_numbers:Vec<String>)->Self{
        // let message = format!("کد تایید : {}",code); 
        let message = format!("{}",code); 
        Self{
            message,
            line_number,
            phone_numbers,
        }
    }
    pub async fn send_sms(&self)->Result<(),Error>{
        let res = reqwest::Client::new()
        .post("https://api.sms.ir/v1/send/likeToLike")
        .header("Content-Type", "application/json")
        .header("X-API-KEY", API_KEY)
        .json(&serde_json::json!({
            "lineNumber":&self.line_number,
            "messageTexts":[&self.message],
            "mobiles":&self.phone_numbers
        }))
        .send()
        .await
        .map_err(|_|return Error::InternalError("sms service #323".to_owned()))?;
        if res.status() != StatusCode::OK{
            return Err(Error::ServiceError(format!("could not send sms {} #323",res.status())))
        }
        Ok(())
    }
}

