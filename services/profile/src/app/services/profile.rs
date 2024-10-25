use sqlx::{Pool, Postgres, Row};

use crate::app::models;
use crate::profile_proto::{AddProfileRequest, Empty, GetProfileBySidRequest, GetProfileByUserIdRequest, GetProfileByUsernameRequest, Profile, UpdateProfileRequest, UpdateUsernameRequest};
use crate::app::types::error::Error;

use super::common;

pub async fn add_profile(pg_db_pool:&Pool<Postgres>,data:AddProfileRequest)->Result<Profile,Error>{
    let result = common::profile_exists_by_user_id(&pg_db_pool, data.user_id).await?;
    if result{
        return Err(Error::ServiceError("profile exists".to_owned()))
    }
    let profile = models::profile::Profile::new(data.user_id);
    let update_log = format!("creating profile {}",chrono::Utc::now().to_string());
    sqlx::query("INSERT INTO profiles (user_id, display_sid, username, display_username, avatar, updated_at) VALUES ($1, $2, $3, $4, 'person', ARRAY[$5])")
    .bind(&profile.user_id)
    .bind(&profile.display_sid)
    .bind(&profile.username)
    .bind(&profile.display_username)
    .bind(update_log)
    .execute(pg_db_pool)
    .await
    .map_err(|err|return Error::InternalError(err.to_string()))?;
    return Ok(Profile { user_id: profile.user_id, display_sid: profile.display_sid, display_username: profile.display_username, username: profile.username,avatar:"person".to_owned() })
}

pub async fn update_username(pg_db_pool:&Pool<Postgres>,data:UpdateUsernameRequest)->Result<Empty,Error>{
    let result = common::profile_exists_by_user_id(&pg_db_pool, data.user_id).await?;
    if !result{
        return Err(Error::NotFoundError("profile not found".to_owned()))
    }
    let username = data.username.to_lowercase();
    let result = common::profile_exists_by_username(&pg_db_pool, &username).await?;
    if result{
        return Err(Error::ServiceError("username already taken".to_owned()))
    }
    let update_log = format!("update username {}",chrono::Utc::now());

    sqlx::query("UPDATE profiles 
    SET 
        username = $1, 
        updated_at = append_array(updated_at, $2) 
    WHERE user_id = $3")
    .bind(username)
    .bind(update_log)
    .bind(data.user_id)
    .execute(pg_db_pool)
    .await
    .map_err(|err|return Error::InternalError(err.to_string()))?;

    return Ok(Empty {})
}


pub async fn update_profile(pg_db_pool:&Pool<Postgres>,data:UpdateProfileRequest)->Result<Empty,Error>{
    let result = common::profile_exists_by_user_id(&pg_db_pool, data.user_id).await?;
    if !result{
        return Err(Error::NotFoundError("profile not found".to_owned()))
    }
    let update_log = format!("update avatar {}",chrono::Utc::now());

    sqlx::query("UPDATE profiles SET 
    updated_at = append_array(updated_at, $1), 
    avatar = COALESCE($2, avatar), 
    display_username = COALESCE($3, display_username) 
    WHERE user_id = $4")
    .bind(update_log)
    .bind(data.avatar)
    .bind(data.display_username)
    .bind(data.user_id)
    .execute(pg_db_pool)
    .await
    .map_err(|err|return Error::InternalError(err.to_string()))?;

    return Ok(Empty {})
}

pub async fn get_profile_by_sid(pg_db_pool:&Pool<Postgres>,data:GetProfileBySidRequest) -> Result<Profile,Error>{
    let row = sqlx::query("SELECT user_id, display_sid, username, display_username, avatar FROM profiles WHERE display_sid = $1") 
    .bind(data.sid)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none(){
        return Err(Error::NotFoundError("user profile not found #404".to_owned()))
    }
    let row = row.unwrap();

    return Ok(Profile { 
        user_id: row.get::<i32,_>("user_id"), 
        display_sid: row.get::<String,_>("display_sid"), 
        display_username: row.get::<String,_>("username"), 
        username:  row.get::<String,_>("username"),
        avatar:row.get::<String,_>("avatar")
    })
}

pub async fn get_profile_by_username(pg_db_pool:&Pool<Postgres>,data:GetProfileByUsernameRequest) -> Result<Profile,Error>{
    let username = data.username.to_lowercase();
    let row = sqlx::query("SELECT user_id, display_sid, username, display_username, avatar FROM profiles WHERE username = $1") 
    .bind(username)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none(){
        return Err(Error::NotFoundError("user profile not found #404".to_owned()))
    }
    let row = row.unwrap();

    return Ok(Profile { 
        user_id: row.get::<i32,_>("user_id"), 
        display_sid: row.get::<String,_>("display_sid"), 
        display_username: row.get::<String,_>("username"), 
        username:  row.get::<String,_>("username"),
        avatar:row.get::<String,_>("avatar")
    })
}

pub async fn get_profile_by_user_id(pg_db_pool:&Pool<Postgres>,data:GetProfileByUserIdRequest) -> Result<Profile,Error>{
    let row = sqlx::query("SELECT user_id, display_sid, username, display_username, avatar FROM profiles WHERE user_id = $1") 
    .bind(data.user_id)
    .fetch_optional(pg_db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none(){
        return Err(Error::NotFoundError("user profile not found #404".to_owned()))
    }
    let row = row.unwrap();

    return Ok(Profile { 
        user_id: row.get::<i32,_>("user_id"), 
        display_sid: row.get::<String,_>("display_sid"), 
        display_username: row.get::<String,_>("username"), 
        username:  row.get::<String,_>("username"),
        avatar:row.get::<String,_>("avatar")
    })
}
