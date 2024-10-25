use sqlx::{Pool, Postgres, Row};
use crate::app::{models, types::error::Error};

pub async fn profile_exists_by_user_id(pg_db_pool:&Pool<Postgres>,user_id:i32)->Result<bool,Error>{
    let row = sqlx::query("SELECT user_id FROM profiles WHERE user_id = $1") 
    .bind(user_id)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_some(){
        return Ok(true) 
    }
    Ok(false)
}

pub async fn profile_exists_by_username(pg_db_pool:&Pool<Postgres>,username:&String)->Result<bool,Error>{
    let row = sqlx::query("SELECT user_id FROM profiles WHERE username = $1") 
    .bind(username)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_some(){
        return Ok(true) 
    }
    Ok(false)
}
pub async fn get_profile_by_user_id(pg_db_pool:&Pool<Postgres>,user_id:i32)->Result<models::profile::Profile,Error>{
    let row = sqlx::query("SELECT user_id, display_sid, username, display_username, avatar FROM profiles WHERE user_id = $1") 
    .bind(user_id)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none(){
        return Err(Error::NotFoundError("user profile not found #404".to_owned()))
    }
    let row = row.unwrap();
    let profile = models::profile::Profile{
        user_id:row.get::<i32,_>("user_id"),
        display_sid:row.get::<String,_>("display_sid"),
        username:row.get::<String,_>("username"),
        display_username:row.get::<String,_>("display_username"),
        avatar:row.get::<String,_>("avatar")
    };
    return Ok(profile)
}