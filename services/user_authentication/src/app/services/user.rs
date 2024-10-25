use chrono::Utc;
use sqlx::{Postgres, Row};

use crate::app::models::user;
use crate::user_proto::{Pagination, Users, User, ChangeUserStatusRequest, Empty, DeleteUserRequest};
use crate::app::types::error::Error;

pub async fn get_users(db_pool:&sqlx::Pool<Postgres>,data:&Pagination)->Result<Users,Error>{
    let rows = sqlx::query("SELECT user_id,phone_number,user_status,created_at FROM users")
    .fetch_all(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;
    
    let mut users = Vec::<User>::new();
    for row in rows{
        let user = User{
            phone_number:row.get::<String,_>("phone_number"),
            user_id:row.get::<i32,_>("user_id"),
            user_status:row.get::<String,_>("user_status"),
            created_at:row.get::<chrono::DateTime<Utc>,_>("created_at").to_string()
        };
        users.push(user);
    }
    if data.get_total{
        let count = sqlx::query("SELECT Count(user_id) AS count FROM users")
        .fetch_one(db_pool)
        .await
        .map_err(|e|return Error::InternalError(e.to_string()))?;
        return Ok(Users { user: users, total_count: Some(count.get::<i64,_>("count")) })
    }
    Ok(Users { user: users, total_count: None })
}

pub async fn change_user_status(db_pool:&sqlx::Pool<Postgres>,data:&ChangeUserStatusRequest)->Result<Empty,Error>{
    let row = sqlx::query("SELECT user_status FROM users WHERE user_id = $1")
    .bind(data.user_id.clone())
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;
    if row.is_none(){
        return Err(Error::NotFoundError("user not found #404".to_owned()))
    }

    let row = row.unwrap();
    let user_status = row.get::<String,_>("user_status");

    let user_status = user::Status::from(user_status);
    let data_user_status = user::Status::from(data.user_status); 

    if user_status == data_user_status{
        return Err(Error::ServiceError(format!("the user status alread is {} #400",data_user_status.to_string())))
    }

    sqlx::query("UPDATE users SET user_status = $1 WHERE user_id = $2")
    .bind(data_user_status.to_string())
    .bind(data.user_id.clone())
    .execute(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    Ok(Empty{})
}

pub async fn delete_user(db_pool:&sqlx::Pool<Postgres>,data:&DeleteUserRequest)->Result<Empty,Error>{
    sqlx::query("DELETE FROM users WHERE user_id = $1")
    .bind(data.user_id)
    .execute(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;
    Ok(Empty{})
}