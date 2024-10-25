use chrono::Utc;
use sqlx::{Postgres, Row};

use crate::account_proto::{
    Empty, GetLoginHistories, GetSessionsRequest, KillSessionRequest, LoginHistories, LoginHistory, LogoutRequest, Pagination, Session, Sessions,
};
use crate::app::types::error::Error;

pub async fn logout(db_pool: &sqlx::Pool<Postgres>, data: LogoutRequest) -> Result<Empty, Error> {
    sqlx::query("DELETE FROM tokens WHERE access_token = $1 AND user_id = $2")
        .bind(data.access_token)
        .bind(data.user_id)
        .execute(db_pool)
        .await
        .map_err(|_| Error::InternalError("failed to logout".to_owned()))?;
    Ok(Empty {})
}

pub async fn kill_session(db_pool: &sqlx::Pool<Postgres>, data: KillSessionRequest) -> Result<Empty, Error> {
    sqlx::query("DELETE FROM tokens WHERE session_id = $1")
        .bind(data.session_id)
        .execute(db_pool)
        .await
        .map_err(|_| Error::InternalError("failed to kill session".to_owned()))?;
    Ok(Empty {})
}

pub async fn get_sessions(db_pool: &sqlx::Pool<Postgres>, data: GetSessionsRequest) -> Result<Sessions, Error> {
    let rows = sqlx::query("SELECT session_id, agent, ip, created_at, token_status FROM tokens WHERE user_id = $1")
        .bind(data.user_id)
        .fetch_all(db_pool)
        .await
        .map_err(|_| Error::InternalError("failed to fetch session".to_owned()))?;

    let sessions: Vec<Session> = rows.into_iter().map(|row| {
        Session {
            session_id: row.get::<i16, _>("session_id") as i32,
            agent: row.get::<String, _>("agent"),
            ip: row.get::<String, _>("ip"),
            created_at: row.get::<chrono::DateTime<Utc>, _>("created_at").to_string(),
            status: row.get::<String, _>("token_status"),
        }
    }).collect();

    Ok(Sessions { sessions })
}

pub async fn get_login_histories(db_pool: &sqlx::Pool<Postgres>, data: GetLoginHistories) -> Result<LoginHistories, Error> {
    let pagination = data.pagination.unwrap_or_else(|| Pagination {
        get_total: false,
        limit: 5,
        offset: 0,
    });

    let rows = sqlx::query(
        "SELECT id, user_id, user_role, section, ip, agent, logged_in_at 
         FROM login_histories 
         WHERE user_id = $1 
         ORDER BY logged_in_at DESC
         OFFSET $2 
         LIMIT $3;"
    )
    .bind(data.user_id)
    .bind(pagination.offset)
    .bind(pagination.limit)
    .fetch_all(db_pool)
    .await
    .map_err(|_| Error::InternalError("failed to fetch login history".to_owned()))?;

    let login_histories: Vec<LoginHistory> = rows.into_iter().map(|row| {
        LoginHistory {
            id: row.get::<String, _>("id"),
            user_id: row.get::<i32, _>("user_id"),
            user_role: row.get::<String, _>("user_role"),
            section: row.get::<String, _>("section"),
            ip: row.get::<String, _>("ip"),
            agent: row.get::<String, _>("agent"),
            logged_in_at: row.get::<chrono::DateTime<Utc>, _>("logged_in_at").to_string(),
        }
    }).collect();

    if pagination.get_total {
        let count = sqlx::query("SELECT COUNT(id) AS count FROM login_histories WHERE user_id = $1")
            .bind(data.user_id)
            .fetch_one(db_pool)
            .await
            .map_err(|_| Error::InternalError("failed to count".to_owned()))?;

        return Ok(LoginHistories {
            login_histories,
            total: Some(count.get::<i64, _>("count")),
        });
    }

    Ok(LoginHistories {
        login_histories,
        total: None,
    })
}
