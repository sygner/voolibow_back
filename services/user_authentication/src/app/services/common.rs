use crate::app::models::user::{Role, Status};
use crate::app::models::{token::Token, user::User};
use crate::app::models::{history_login, token};
use crate::app::types::error::Error;
use bb8_redis::RedisConnectionManager;
use bb8_redis::redis::AsyncCommands;
use chrono::Utc;
use sqlx::{Postgres, Row};

pub async fn user_exists_by_phone_pg(db_pool:&sqlx::Pool<Postgres>,phone_number:&String)->Result<(),Error>{
    let data = sqlx::query("SELECT user_id FROM users WHERE phone_number = $1").bind(phone_number).fetch_optional(db_pool)
    .await
    .map_err(|e| Error::InternalError(e.to_string()))?;

    if data.is_some(){
        return Err(Error::ServiceError(format!("phone number already exists")))
    }
    Ok(())
}


pub async fn phone_number_exists_rd(db_pool:&bb8::Pool<RedisConnectionManager>,phone_number:&String)->Result<(),Error>{
    let mut rd_db_pool =  db_pool.get()
    .await
    .map_err(|_|return Error::InternalError("try later #555".to_owned()))?;

    let exists:Option<String> = rd_db_pool.get(phone_number).await.map_err(|e|return Error::InternalError(e.to_string()))?;
    if exists.is_none(){
        return Err(Error::NotFoundError(String::new()))
    }

    Ok(())
}

pub async fn create_token(db_pool:&sqlx::Pool<Postgres>,token_data:&Token)->Result<(),Error>{
    sqlx::query("INSERT INTO tokens (access_token,refresh_token,user_id,user_role,session_id,token_status,ip,agent,created_at,access_token_expire_at,refresh_token_expire_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)")
    .bind(token_data.access_token.clone())
    .bind(token_data.refresh_token.clone())
    .bind(token_data.user_id)
    .bind(token_data.role.to_string())
    .bind(token_data.session_id)
    .bind(token_data.status.to_string())
    .bind(token_data.ip.clone())
    .bind(token_data.agent.clone())
    .bind(token_data.created_at)
    .bind(token_data.access_token_expire_at)
    .bind(token_data.refresh_token_expire_at)
    .execute(db_pool)
    .await
    .map_err(|_|return Error::InternalError("failed to create".to_owned()))?;
    return Ok(())
}

// pub async fn get_user_id_by_phone_number_pg(db_pool:&sqlx::Pool<Postgres>,phone_number:&String)->Result<Option<i32>,Error>{
//     let row = sqlx::query("SELECT user_id FROM users WHERE phone_number = $1")
//     .bind(phone_number)
//     .fetch_optional(db_pool)
//     .await
//     .map_err(|e| Error::InternalError(e.to_string()))?;
//     if row.is_none(){
//         return Ok(None)
//     }
//     let row = row.unwrap();
//     let user_id = row.get::<i32,_>("user_id");
//     Ok(Some(user_id))
// }

pub async fn get_user_by_phone_number_pg(db_pool: &sqlx::Pool<Postgres>, phone_number: &String) -> Result<Option<User>, Error> {
    let row = sqlx::query("SELECT user_id, phone_number, user_role, user_status, created_at FROM users WHERE phone_number = $1")
        .bind(phone_number)
        .fetch_optional(db_pool)
        .await
        .map_err(|_| Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none() {
        return Ok(None);
    }

    let row = row.unwrap();

    let user = User {
        user_id: row.get::<i32, _>("user_id"),
        phone_number: row.get::<String, _>("phone_number"),
        role: Role::from(row.get::<String, _>("user_role")),
        user_status: Status::from(row.get::<String, _>("user_status")),
        created_at: row.get::<chrono::DateTime<Utc>, _>("created_at"),
    };

    Ok(Some(user))
}


pub async fn get_token_by_access_token_and_agent(db_pool:&sqlx::Pool<Postgres>,access_token:&String,agent:&String,role:&String)->Result<Token,Error>{
    let row = sqlx::query("SELECT refresh_token,user_id,session_id,token_status,ip,agent,user_role,created_at,access_token_expire_at,refresh_token_expire_at FROM tokens WHERE access_token = $1 AND agent = $2 AND user_role = $3")
    .bind(access_token)
    .bind(agent)
    .bind(role)
    .fetch_optional(db_pool)
    .await
    .map_err(|_| Error::InternalError("failed to fetch".to_owned()))?;

    if row.is_none(){
        return Err(Error::NotFoundError("token not found #404".to_owned()))
    }

    let row = row.unwrap();
    // let refresh_token = row.get::<String,_>("refresh_token");
    // let user_id = row.get::<i32,_>("user_id");
    // let session_id = row.get::<i32,_>("session_id");
    // let ip = row.get::<String,_>("ip");
    // let agent = row.get::<String,_>("agent");
    // let created_at = row.get::<chrono::DateTime<Local>,_>("created_at");
    // let access_token_expiry = row.get::<chrono::DateTime<Local>,_>("access_token_expire_at");
    // let refresh_token_expiry = row.get::<chrono::DateTime<Local>,_>("refresh_token_expire_at");
    // NaiveDateTime::from_timestamp_millis(timestamp).map(|e|Error::InternalError("timestamp error".to_owned()))?;
    let token = Token{
        user_id:row.get::<i32,_>("user_id"),
        access_token:access_token.to_owned(),
        refresh_token:row.get::<String,_>("refresh_token"),
        session_id:row.get::<i16,_>("session_id"),
        role:row.get::<String,_>("user_role"),
        agent:row.get::<String,_>("agent"),
        ip:row.get::<String,_>("ip"),
        status:token::Status::from(row.get::<String,_>("token_status")),
        created_at: row.get::<chrono::DateTime<Utc>,_>("created_at"),
        access_token_expire_at: row.get::<chrono::DateTime<Utc>,_>("access_token_expire_at"),
        refresh_token_expire_at: row.get::<chrono::DateTime<Utc>,_>("refresh_token_expire_at"),
    };
    return Ok(token)
}


// pub async fn get_token_by_access_token(db_pool:&sqlx::Pool<Postgres>,access_token:&String)->Result<Token,Error>{
//     let row = sqlx::query("SELECT (refresh_token,user_id,session_id,token_status,ip,agent,created_at,access_token_expire_at,refresh_token_expire_at) FROM tokens WHERE access_token = $1")
//     .bind(access_token)
//     .fetch_optional(db_pool)
//     .await
//     .map_err(|e| Error::InternalError(e.to_string()))?;

//     if row.is_none(){
//         return Err(Error::NotFoundError("token not found #404".to_owned()))
//     }

//     let row = row.unwrap();
//     let status = row.get::<String,_>("token_status");
   
//     let status = token::Status::from(status);
//     let token = Token{
//         user_id:row.get::<i32,_>("user_id"),
//         access_token:access_token.to_owned(),
//         refresh_token:row.get::<String,_>("refresh_token"),
//         role:row.get::<String,_>("user_role"),
//         session_id:row.get::<i16,_>("session_id"),
//         agent:row.get::<String,_>("agent"),
//         ip:row.get::<String,_>("ip"),
//         status:status,
//         created_at:row.get::<chrono::DateTime<Utc>,_>("created_at"),
//         access_token_expire_at: row.get::<chrono::DateTime<Utc>,_>("access_token_expire_at"),
//         refresh_token_expire_at: row.get::<chrono::DateTime<Utc>,_>("refresh_token_expire_at"),
//     };
//     return Ok(token)
// }


pub async fn delete_token_by_access_token_and_refresh_token(db_pool:&sqlx::Pool<Postgres>,access_token:&String,refresh_token:&String)->Result<(),Error>{
    sqlx::query("DELETE FROM tokens WHERE access_token = $1 AND refresh_token = $2")
    .bind(access_token)
    .bind(refresh_token)
    .execute(db_pool)
    .await
    .map_err(|_| Error::InternalError("failed to delete".to_owned()))?;
    return Ok(())
}


// pub async fn delete_token_by_access_token(db_pool:&sqlx::Pool<Postgres>,access_token:&String)->Result<(),Error>{
//     sqlx::query("DELETE FROM tokens WHERE access_token = $1")
//     .bind(access_token)
//     .execute(db_pool)
//     .await
//     .map_err(|e| Error::InternalError(e.to_string()))?;
//     return Ok(())
// }

// pub async fn get_user_role_by_user_id(db_pool:&sqlx::Pool<Postgres>,user_id:&i32) ->Result<Option<String>,Error>{
//     let row = sqlx::query("SELECT user_role FROM users WHERE user_id = $1")
//     .bind(user_id)
//     .fetch_optional(db_pool)
//     .await
//     .map_err(|e| Error::InternalError(e.to_string()))?;
//     if row.is_none(){
//         return Ok(None)
//     }else{
//         let row = row.unwrap();
//         return Ok(Some(row.get::<String,_>("user_role")))
//     }
// }


pub async fn add_history_login(db_pool:&sqlx::Pool<Postgres>,data:history_login::LoginHistory)->Result<(),Error>{
    sqlx::query("INSERT INTO login_histories (id, user_id, user_role, section, ip, agent, logged_in_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())")
    .bind(data.id)
    .bind(data.user_id)
    .bind(data.user_role)
    .bind(data.section)
    .bind(data.ip)
    .bind(data.agent)
    .execute(db_pool)
    .await
    .map_err(|_|Error::InternalError("failed to add".to_owned()))?;
    return Ok(())
}