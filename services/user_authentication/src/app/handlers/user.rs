use std::sync::Arc;

use sqlx::Postgres;
use tonic::{Request, Response, Status};
use crate::app::services;
use crate::user_proto::user_service_server::UserService;
use crate::user_proto::{ChangeUserStatusRequest,Empty,Pagination, Users, DeleteUserRequest};


pub struct UserHandler{
    pub postgres_db:Arc<sqlx::Pool<Postgres>>,
}

impl UserHandler{
    pub fn new(postgres_db:Arc<sqlx::Pool<Postgres>>)->UserHandler{
        Self{
            postgres_db
        }
    }
}

#[tonic::async_trait]
impl UserService for UserHandler{
    async fn get_users(&self,request:Request<Pagination>)->Result<Response<Users>,Status>{
        let res = services::user::get_users(&self.postgres_db, &request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn change_user_status(&self,request:Request<ChangeUserStatusRequest>)->Result<Response<Empty>,Status>{
        let res = services::user::change_user_status(&self.postgres_db, &request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn delete_user(&self,request:Request<DeleteUserRequest>)->Result<Response<Empty>,Status>{
        let res = services::user::delete_user(&self.postgres_db, &request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
}