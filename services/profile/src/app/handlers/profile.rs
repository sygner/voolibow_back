use std::sync::Arc;

use sqlx::Postgres;
use tonic::{Request, Response, Status};

use crate::app::services;
use crate::profile_proto::{AddProfileRequest, Empty, GetProfileBySidRequest, GetProfileByUserIdRequest, GetProfileByUsernameRequest, Profile, UpdateProfileRequest, UpdateUsernameRequest};
use crate::profile_proto::profile_service_server::ProfileService;
pub struct ProfileHandler{
    pub postgres_db:Arc<sqlx::Pool<Postgres>>,
}

impl ProfileHandler{
    pub fn new(postgres_db:Arc<sqlx::Pool<Postgres>>)->Self{
        Self { postgres_db }
    }
}

#[tonic::async_trait]
impl ProfileService for ProfileHandler{
    async fn add_profile(&self,request:Request<AddProfileRequest>)->Result<Response<Profile>,Status>{
        let res = services::profile::add_profile(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn update_username(&self,request:Request<UpdateUsernameRequest>)->Result<Response<Empty>,Status>{
        let res = services::profile::update_username(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn update_avatar(&self,request:Request<UpdateProfileRequest>)->Result<Response<Empty>,Status>{
        let res = services::profile::update_profile(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_profile_by_sid(&self,request:Request<GetProfileBySidRequest>)->Result<Response<Profile>,Status>{
        let res = services::profile::get_profile_by_sid(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_profile_by_username(&self,request:Request<GetProfileByUsernameRequest>)->Result<Response<Profile>,Status>{
        let res = services::profile::get_profile_by_username(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }

    async fn get_profile_by_user_id(&self,request:Request<GetProfileByUserIdRequest>)->Result<Response<Profile>,Status>{
        let res = services::profile::get_profile_by_user_id(&self.postgres_db.as_ref(), request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
}