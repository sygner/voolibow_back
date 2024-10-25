use serde::Serialize;
use tonic::Status;


#[derive(Debug,Serialize,thiserror::Error,PartialEq, Eq)]
pub enum Error{
    #[error("Internal Error: `{0}`")]
    InternalError(String),
    #[error("Service Error: `{0}`")]
    ServiceError(String),
    #[error("Not Found Error: `{0}`")]
    NotFoundError(String),
    #[error("Permission Denied Error: `{0}`")]
    PermissionDeniedError(String),
    #[error("Database Error: `{0}`")]
    DatabaseError(String),
    #[error("Pool Error: `{0}`")]
    DBPoolError(String),
}

impl Error{
    pub fn to_status(&self)->Status{
        match self{
            Error::InternalError(m)=>Status::internal(m),
            Error::ServiceError(m)=>Status::aborted(m),
            Error::NotFoundError(m)=>Status::not_found(m),
            Error::PermissionDeniedError(m)=>Status::permission_denied(m),
            _=>Status::unknown("unkown error")
        }
    }
}