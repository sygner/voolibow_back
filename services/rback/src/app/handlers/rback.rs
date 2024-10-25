use std::sync::Arc;

use sqlx::Postgres;
use tonic::{Request, Status, Response};
use crate::app::services;
use crate::rback_proto::r_back_service_server::RBackService;
use crate::rback_proto::{
    Empty,
    Roles,
    Permissions, 
    VerifyUserPermissionRequest, 
    GetUserRolesRequest, 
    GetRolePermissionsRequest, AddUserRoleRequest, AddRolePermissionRequest, RolePermissions,
    
    };
pub struct RBackHandler{
    pub postgres_db:Arc<sqlx::Pool<Postgres>>,
}

impl RBackHandler{
    pub fn new(postgres_db:Arc<sqlx::Pool<Postgres>>)->Self{
        Self { postgres_db }
    }
}

#[tonic::async_trait]
impl RBackService for RBackHandler{
    async fn verify_user_permission(&self,request:Request<VerifyUserPermissionRequest>)->Result<Response<Empty>,Status>{
        let res = services::rback::verify_user_permission(&self.postgres_db, request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_all_roles(&self,_:Request<Empty>)->Result<Response<Roles>,Status>{
        let res = services::rback::get_all_roles(&self.postgres_db).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_all_permissions(&self,_:Request<Empty>)->Result<Response<Permissions>,Status>{
        let res = services::rback::get_all_permissions(&self.postgres_db).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_role_permissions(&self,request:Request<GetRolePermissionsRequest>)->Result<Response<Permissions>,Status>{
        let res = services::rback::get_role_permissions(&self.postgres_db,request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn get_user_roles(&self,request:Request<GetUserRolesRequest>)->Result<Response<RolePermissions>,Status>{
        let res = services::rback::get_user_permissions(&self.postgres_db,request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn add_user_role(&self,request:Request<AddUserRoleRequest>)->Result<Response<Empty>,Status>{
        let res = services::rback::add_user_role(&self.postgres_db,request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
    async fn add_role_permission(&self,request:Request<AddRolePermissionRequest>)->Result<Response<Empty>,Status>{
        let res = services::rback::add_role_permission(&self.postgres_db,request.into_inner()).await.map_err(|e| return e.to_status())?;
        Ok(Response::new(res))
    }
}