use std::sync::Arc;

use sqlx::Postgres;
use tonic::{Request, Response, Status};
use crate::account_proto::{Empty, GetLoginHistories, GetSessionsRequest, KillSessionRequest, LoginHistories, LogoutRequest, Sessions};
use crate::account_proto::account_service_server::AccountService;
use crate::app::services::account;

pub struct AccountHandler{
    pub postgres_db:Arc<sqlx::Pool<Postgres>>,
}

impl AccountHandler{
    pub fn new(postgres_db:Arc<sqlx::Pool<Postgres>>)->Self{
        Self { postgres_db }
    }
}

#[tonic::async_trait]
impl AccountService for AccountHandler{
    async fn logout(&self,request:Request<LogoutRequest>)->Result<Response<Empty>,Status>{
        let res = account::logout(&self.postgres_db, request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn kill_session(&self,request:Request<KillSessionRequest>)->Result<Response<Empty>,Status>{
        let res = account::kill_session(&self.postgres_db, request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_sessions(&self,request:Request<GetSessionsRequest>)->Result<Response<Sessions>,Status>{
        let res = account::get_sessions(&self.postgres_db, request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_login_history(&self,request:Request<GetLoginHistories>)->Result<Response<LoginHistories>,Status>{
        let res = account::get_login_histories(&self.postgres_db, request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
}