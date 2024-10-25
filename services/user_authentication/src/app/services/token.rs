use chrono::Utc;
use sqlx::{Postgres, Row};
use crate::app::models::token;
use crate::token_proto::{VerificationRequest, VerificationResponse, RenewTokenRequest, TokenInfo, ChangeTokenStatusRequest, Empty, Pagination, Tokens, Token};
use crate::app::types::error::Error;
use super::common::{self, delete_token_by_access_token_and_refresh_token};

pub async fn verify_token(db_pool: &sqlx::Pool<Postgres>, data: VerificationRequest) -> Result<VerificationResponse, Error> {
    let token = common::get_token_by_access_token_and_agent(db_pool, &data.access_token, &data.agent, &data.role).await?;
    if !token.validate_expiry() || !token.validate_status() {
        return Err(Error::PermissionDeniedError("Token is invalid or expired #777".to_owned()));
    }
    Ok(VerificationResponse {
        user_id: token.user_id,
        session_id: token.session_id as i32,
        role: token.role,
    })
}

pub async fn renew_token(db_pool: &sqlx::Pool<Postgres>, data: RenewTokenRequest, expiry: i32) -> Result<TokenInfo, Error> {
    let row = sqlx::query("SELECT user_id, user_role FROM tokens WHERE access_token = $1 AND refresh_token = $2 AND agent = $3")
        .bind(data.access_token.clone())
        .bind(data.refresh_token.clone())
        .bind(data.agent.clone())
        .fetch_optional(db_pool)
        .await
        .map_err(|e| Error::InternalError(e.to_string()))?;

    if let Some(row) = row {
        delete_token_by_access_token_and_refresh_token(db_pool, &data.access_token, &data.refresh_token).await?;
        let user_id = row.get::<i32, _>("user_id");
        let user_role = row.get::<String, _>("user_role");
        let new_token = token::Token::new(user_id, data.agent, data.ip, user_role, expiry);
        common::create_token(db_pool, &new_token).await?;
        Ok(TokenInfo {
            access_token: new_token.access_token,
            refresh_token: new_token.refresh_token,
            expiry,
        })
    } else {
        Err(Error::NotFoundError("Token not found #404".to_owned()))
    }
}

pub async fn change_token_status(db_pool: &sqlx::Pool<Postgres>, data: ChangeTokenStatusRequest) -> Result<Empty, Error> {
    let row = sqlx::query("SELECT token_status FROM tokens WHERE access_token = $1")
        .bind(data.access_token.clone())
        .fetch_optional(db_pool)
        .await
        .map_err(|e| Error::InternalError(e.to_string()))?;

    if let Some(row) = row {
        let token_status = token::Status::from(row.get::<String, _>("token_status"));
        let data_token_status = token::Status::from(data.token_status);

        if token_status == data_token_status {
            return Err(Error::ServiceError(format!("The token status already is {} #400", data_token_status.to_string())));
        }

        sqlx::query("UPDATE tokens SET token_status = $1 WHERE access_token = $2")
            .bind(data_token_status.to_string())
            .bind(data.access_token.clone())
            .execute(db_pool)
            .await
            .map_err(|e| Error::InternalError(e.to_string()))?;
        Ok(Empty {})
    } else {
        Err(Error::NotFoundError("Token not found #404".to_owned()))
    }
}

pub async fn get_tokens(db_pool: &sqlx::Pool<Postgres>, data: Pagination) -> Result<Tokens, Error> {
    let rows = sqlx::query(
        "SELECT access_token, refresh_token, user_id, session_id, token_status, ip, agent, created_at, access_token_expire_at, refresh_token_expire_at 
         FROM tokens 
         ORDER BY created_at DESC
         OFFSET $1 
         LIMIT $2;"
    )
    .bind(data.offset)
    .bind(data.limit)
    .fetch_all(db_pool)
    .await
    .map_err(|_| Error::InternalError("failed to fetch".to_owned()))?;

    let tokens: Vec<Token> = rows.into_iter().map(|row| {
        Token {
            access_token: row.get::<String, _>("access_token"),
            refresh_token: row.get::<String, _>("refresh_token"),
            user_id: row.get::<i32, _>("user_id"),
            session_id: row.get::<i16, _>("session_id") as i32,
            token_status: row.get::<String, _>("token_status"),
            ip: row.get::<String, _>("ip"),
            agent: row.get::<String, _>("agent"),
            created_at: row.get::<chrono::DateTime<Utc>, _>("created_at").to_string(),
            access_token_expire_at: row.get::<chrono::DateTime<Utc>, _>("access_token_expire_at").to_string(),
            refresh_token_expire_at: row.get::<chrono::DateTime<Utc>, _>("refresh_token_expire_at").to_string(),
            role: row.get::<String, _>("user_role"),
        }
    }).collect();

    if data.get_total {
        let count = sqlx::query("SELECT COUNT(access_token) AS count FROM tokens")
            .fetch_one(db_pool)
            .await
            .map_err(|_| Error::InternalError("failed to count".to_owned()))?;
        return Ok(Tokens {
            tokens,
            total_count: Some(count.get::<i64, _>("count")),
        });
    }

    Ok(Tokens {
        tokens,
        total_count: None,
    })
}
