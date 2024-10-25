use std::sync::Arc;

use sqlx::Postgres;
use tonic::{Request, Response, Status};
use crate::app::services::token;
use crate::token_proto::token_service_server::TokenService;
use crate::token_proto::{
    VerificationRequest,
    VerificationResponse,
    Empty,
    RenewTokenRequest,
    TokenInfo,
    ChangeTokenStatusRequest, Tokens, Pagination};
pub struct TokenHandler{
    pub postgres_db:Arc<sqlx::Pool<Postgres>>,
    pub token_life_expiry:i32
}


impl TokenHandler{
    pub fn new(postgres_db:Arc<sqlx::Pool<Postgres>>,token_life_expiry:i32) ->Self{
        Self { 
            postgres_db,
            token_life_expiry
        }
    }
}

#[tonic::async_trait]
impl TokenService for TokenHandler{
    async fn verify_token(&self,request:Request<VerificationRequest>)->Result<Response<VerificationResponse>,Status>{
        let res = token::verify_token(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn renew_token(&self,request:Request<RenewTokenRequest>)->Result<Response<TokenInfo>,Status>{
        let res = token::renew_token(&self.postgres_db.as_ref(), request.into_inner(), self.token_life_expiry).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn change_token_status(&self,request:Request<ChangeTokenStatusRequest>)->Result<Response<Empty>,Status>{
        let res = token::change_token_status(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_tokens(&self,request:Request<Pagination>)->Result<Response<Tokens>,Status>{
        let res = token::get_tokens(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
}
