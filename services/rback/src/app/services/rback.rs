use sqlx::{Postgres, Row};

use crate::app::types::error::Error;
use crate::rback_proto::{Empty, VerifyUserPermissionRequest, Roles, Role, Permissions, Permission, GetRolePermissionsRequest, GetUserRolesRequest, RolePermissions, RolePermission, AddUserRoleRequest, AddRolePermissionRequest};
  
pub async fn verify_user_permission(db_pool:&sqlx::Pool<Postgres>,data:VerifyUserPermissionRequest)->Result<Empty,Error>{
    let row = sqlx::query("
    SELECT user_id FROM user_roles
	INNER JOIN roles
	ON user_roles.role_id=roles.role_id
	INNER JOIN permissions
	ON roles.permission_id = permissions.permission_id
	WHERE
	user_roles.user_id=$1 AND
	permissions.permission_id = $2
    ")
    .bind(data.user_id)
    .bind(data.permission_id)
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    if row.is_none(){
        return Err(Error::PermissionDeniedError("no permission #701".to_owned()))
    }
    return Ok(Empty {})
}


pub async fn get_all_roles(db_pool:&sqlx::Pool<Postgres>)->Result<Roles,Error>{
    let rows = sqlx::query("SELECT role_id FROM roles")
    .fetch_all(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;


    let mut roles = Vec::<Role>::new();

    for row in rows{
        let role = Role{
            role_id:row.get::<String,_>("role_id"),
        };
        roles.push(role)
    }

    Ok(Roles { roles })
}

pub async fn get_all_permissions(db_pool:&sqlx::Pool<Postgres>)->Result<Permissions,Error>{
    let rows = sqlx::query("SELECT permission_id,description FROM permissions")
    .fetch_all(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;


    let mut permissions = Vec::<Permission>::new();

    for row in rows{
        let permission = Permission{
            permission_id:row.get::<String,_>("permission_id"),
            description:row.get::<String,_>("description"),
        };
        permissions.push(permission)
    }

    Ok(Permissions { permissions })
}

pub async fn get_role_permissions(db_pool:&sqlx::Pool<Postgres>,data:GetRolePermissionsRequest)->Result<Permissions,Error>{
    let rows = sqlx::query("SELECT permissions.permission_id,permissions.description FROM roles
	LEFT JOIN role_permissions
	ON roles.role_id = $1
	INNER JOIN permissions
	ON role_permissions.permission_id = permissions.permission_id")
    .bind(data.role_id)
    .fetch_all(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;


    
    let mut permissions = Vec::<Permission>::new();

    for row in rows{
        let permission = Permission{
            permission_id:row.get::<String,_>("permission_id"),
            description:row.get::<String,_>("description"),
        };
        permissions.push(permission)
    }

    Ok(Permissions { permissions })
}

pub async fn get_user_permissions(db_pool:&sqlx::Pool<Postgres>,data:GetUserRolesRequest)->Result<RolePermissions,Error>{
    let rows = sqlx::query("SELECT roles.role_id, permissions.permission_id,permissions.description FROM user_roles
	INNER JOIN roles
	ON user_roles.role_id = roles.role_id
	INNER JOIN role_permissions
	ON roles.role_id = role_permissions.role_id
	INNER JOIN permissions
	ON role_permissions.permission_id = permissions.permission_id
	WHERE user_roles.user_id=$1")
    .bind(data.user_id)
    .fetch_all(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;


    
    let mut role_permissions_map = std::collections::HashMap::<String, RolePermission>::new();

    for row in rows {
        let role_id = row.get::<String, _>("role_id");
        let permission = Permission {
            permission_id: row.get::<String, _>("permission_id"),
            description: row.get::<String, _>("description"),
        };

        let role_permission = role_permissions_map
            .entry(role_id.clone())
            .or_insert_with(|| RolePermission {
                role_id: role_id.clone(),
                permissions: Vec::new(),
            });

        role_permission.permissions.push(permission);
    }

    let role_permissions = role_permissions_map.into_values().collect();

    Ok(RolePermissions { role_permissions })
}

pub async fn add_user_role(db_pool:&sqlx::Pool<Postgres>,data:AddUserRoleRequest)->Result<Empty,Error>{
    let row_exists = sqlx::query("SELECT 1 FROM roles WHERE role_id=$1")
    .bind(data.role_id.clone())
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    if row_exists.is_none(){
        return Err(Error::NotFoundError("role id not found #404".to_owned()))
    }

    let row_exists = sqlx::query("SELECT 1 FROM user_roles WHERE user_id=$1")
    .bind(data.user_id.clone())
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    if row_exists.is_some(){
        return Err(Error::NotFoundError("already exists #400".to_owned()))
    }

    sqlx::query("INSERT INTO user_roles VALUES ($1,$2)")
    .bind(data.user_id)
    .bind(data.role_id)
    .execute(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    Ok(Empty{})
}

pub async fn add_role_permission(db_pool:&sqlx::Pool<Postgres>,data:AddRolePermissionRequest)->Result<Empty,Error>{
    let row_exists = sqlx::query("SELECT 1 FROM permissions WHERE permission_id=$1")
    .bind(data.permission_id.clone())
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    if row_exists.is_none(){
        return Err(Error::NotFoundError("permission id not found #404".to_owned()))
    }

    let row_exists = sqlx::query("SELECT 1 FROM roles WHERE role_id=$1")
    .bind(data.role_id.clone())
    .fetch_optional(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    if row_exists.is_none(){
        return Err(Error::NotFoundError("role id not found #404".to_owned()))
    }


    sqlx::query("INSERT INTO role_permissions VALUES ($1,$2)")
    .bind(data.role_id)
    .bind(data.permission_id)
    .execute(db_pool)
    .await
    .map_err(|e|return Error::InternalError(e.to_string()))?;

    Ok(Empty{})
}
