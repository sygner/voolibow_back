use crate::app::models;
use crate::app::models::{history_login::LoginHistory, token::Token, user::User};
use crate::app::types::error::Error;
use crate::authentication_proto::{
    OptionalResponse, SigninRequest, SignupRequest, TokenInfo, VerificationRequest, VerificationResponse,
};
use std::sync::Arc;
use bb8_redis::{redis::AsyncCommands, RedisConnectionManager};
use sqlx::{Postgres, Row};
use super::common;

pub async fn signup(
    pg_db_pool: &sqlx::Pool<Postgres>,
    rd_db_pool: &bb8::Pool<RedisConnectionManager>,
    data: SignupRequest,
) -> Result<OptionalResponse, Error> {
    let rgx = regex::Regex::new(r#"^(\+98|0)?9\d{9}$"#)
        .map_err(|_| Error::InternalError("failed to validate".to_owned()))?;

    if !rgx.is_match(&data.phone) {
        return Err(Error::ServiceError("the format of phone number is incorrect #111".to_owned()));
    }

    match common::phone_number_exists_rd(rd_db_pool, &data.phone).await {
        Ok(_) => return Err(Error::ServiceError("you've already made a process line. try later #111".to_owned())),
        Err(Error::NotFoundError(_)) => {},
        Err(e) => return Err(Error::InternalError(e.to_string())),
    }

    let phone_number = Arc::new(data.phone.clone());
    common::user_exists_by_phone_pg(pg_db_pool, &phone_number).await?;

    let value = serde_json::to_string(&User::new(&phone_number))
        .map_err(|_| Error::InternalError("try later #711".to_owned()))?;

    let mut rd_conn = rd_db_pool.get().await.map_err(|_| Error::InternalError("try later #555".to_owned()))?;
    let code = idgen::numeric_code_i32(10151, 99592);

    // pkg::SMS::new_verification_message(code, "_".to_owned(), vec![phone_number.as_ref().to_owned()]).send_sms().await?;

    rd_conn.set_ex(phone_number.as_ref(), value, 120)
        .await
        .map_err(|_| Error::InternalError("try later #712".to_owned()))?;
    rd_conn.set_ex(code.to_string(), phone_number.as_ref(), 120)
        .await
        .map_err(|_| Error::InternalError("try later #712".to_owned()))?;

    println!("{}", code);

    Ok(OptionalResponse {
        msg: None,
        code: Some(code.to_string()),
    })
}

pub async fn verify(
    pg_db_pool: &sqlx::Pool<Postgres>,
    rd_db_pool: &bb8::Pool<RedisConnectionManager>,
    data: VerificationRequest,
    expiry: i32,
) -> Result<VerificationResponse, Error> {
    let mut rd_conn = rd_db_pool.get().await.map_err(|_| Error::InternalError("try later #555".to_owned()))?;
    
    match data.verification_method {
        0 => {
            let phone_number: Option<String> = rd_conn.get_del(&data.code).await.map_err(|_| Error::InternalError("try later #712".to_owned()))?;
            let phone_number = phone_number.ok_or_else(|| Error::NotFoundError("code not found #404".to_owned()))?;
            let user_json_data: String = rd_conn.get_del(&phone_number).await.map_err(|_| Error::InternalError("try later #712".to_owned()))?;

            let mut user: User = serde_json::from_str(&user_json_data).map_err(|_| Error::InternalError("try later #711".to_owned()))?;
            let user_row = sqlx::query("INSERT INTO users (phone_number, user_role, user_status, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id")
                .bind(&user.phone_number)
                .bind(&user.role.to_string())
                .bind(&user.user_status.to_string())
                .bind(&user.created_at)
                .fetch_one(pg_db_pool)
                .await
                .map_err(|e| Error::InternalError(e.to_string()))?;
            user.user_id = user_row.get("user_id");

            println!("User ID: {}", user.user_id);
            println!("User Role: {}", user.role.to_string());

            let token = Token::new(user.user_id, data.agent.clone(), data.ip.clone(), user.role.to_string(), expiry);
            common::create_token(pg_db_pool, &token).await?;
            
            let login_history = LoginHistory::new(
                user.user_id,
                user.role.to_string(),
                "voolibow".to_owned(),
                data.ip,
                data.agent,
                chrono::Utc::now(),
            );
            common::add_history_login(pg_db_pool, login_history).await?;

            Ok(VerificationResponse {
                token_info: Some(TokenInfo {
                    access_token: token.access_token,
                    refresh_token: token.refresh_token,
                    expiry,
                }),
                user_id: user.user_id,
            })
        }
        1 => {
            let phone_number: Option<String> = rd_conn.get_del(&data.code).await.map_err(|_| Error::InternalError("try later #712".to_owned()))?;
            let phone_number = phone_number.ok_or_else(|| Error::NotFoundError("code not found #404".to_owned()))?;

            match common::user_exists_by_phone_pg(pg_db_pool, &phone_number).await {
                Ok(_) => return Err(Error::NotFoundError("phone number not found #404".to_owned())),
                Err(Error::ServiceError(_)) => {},
                Err(e) => return Err(Error::InternalError(e.to_string())),
            }

            let user_data = common::get_user_by_phone_number_pg(pg_db_pool, &phone_number).await?
                .ok_or_else(|| Error::NotFoundError("phone number not found #404".to_owned()))?;

            let token = Token::new(user_data.user_id, data.agent.clone(), data.ip.clone(), user_data.role.to_string(), expiry);
            common::create_token(pg_db_pool, &token).await?;
            
            let login_history = LoginHistory::new(
                user_data.user_id,
                user_data.role.to_string(),
                "voolibow".to_owned(),
                data.ip,
                data.agent,
                chrono::Utc::now(),
            );
            common::add_history_login(pg_db_pool, login_history).await?;

            Ok(VerificationResponse {
                token_info: Some(TokenInfo {
                    access_token: token.access_token,
                    refresh_token: token.refresh_token,
                    expiry,
                }),
                user_id: user_data.user_id,
            })
        }
        _ => Err(Error::ServiceError("unknown method #400".to_owned())),
    }
}

pub async fn signin(
    pg_db_pool: &sqlx::Pool<Postgres>,
    rd_db_pool: &bb8::Pool<RedisConnectionManager>,
    data: SigninRequest,
) -> Result<OptionalResponse, Error> {
    let phone_number = Arc::new(data.phone.clone());

    match common::user_exists_by_phone_pg(pg_db_pool, &phone_number).await {
        Ok(_) => return Err(Error::NotFoundError("phone number not found #404".to_owned())),
        Err(Error::ServiceError(_)) => {},
        Err(e) => return Err(Error::InternalError(e.to_string())),
    }

    let row = sqlx::query("SELECT user_status FROM users WHERE phone_number = $1")
        .bind(phone_number.as_ref())
        .fetch_one(pg_db_pool)
        .await
        .map_err(|e| Error::InternalError(e.to_string()))?;

    let user_status = models::user::Status::from(row.get::<String, _>("user_status"));
    if !user_status.validate_status() {
        return Err(Error::ServiceError(format!("your account status: {} #400", user_status.to_string())));
    }

    let mut rd_conn = rd_db_pool.get().await.map_err(|_| Error::InternalError("try later #555".to_owned()))?;
    let code = idgen::numeric_code_i32(10151, 99592);

    println!("{}", code);
    // pkg::SMS::new_verification_message(code, "_".to_owned(), vec![phone_number.as_ref().to_owned()]).send_sms().await?;

    rd_conn.set_ex(code.to_string(), phone_number.as_ref(), 120)
        .await
        .map_err(|_| Error::InternalError("try later #712".to_owned()))?;

    Ok(OptionalResponse {
        msg: None,
        code: Some(code.to_string()),
    })
}
