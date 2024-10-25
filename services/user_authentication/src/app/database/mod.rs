use bb8_redis::RedisConnectionManager;
use sqlx::{postgres::PgPoolOptions, Postgres};
use super::types::error::Error;
pub async fn postgres_connection(username:String,password:String,host:String,port:usize,table:String)->Result<sqlx::Pool<Postgres>,Error>{
    let uri = format!("postgres://{}:{}@{}:{}/{}",username,password,host,port,table);
    let pool = PgPoolOptions::new()
    .max_connections(5)
    .connect(uri.as_str()).await.map_err(|e|return Error::DBPoolError(e.to_string()));
    Ok(pool.unwrap())
}


pub async fn redis_connection(host:String)->Result<bb8::Pool<RedisConnectionManager>,Error>{
    let manager = bb8_redis::RedisConnectionManager::new(format!("redis://{}",host)).map_err(|e|return Error::DatabaseError(e.to_string())).unwrap();
    let pool = bb8::Pool::builder()
        .build(manager)
        .await
        .map_err(|e|return Error::DBPoolError(e.to_string()));
    // let res = redis::Client::open(format!("redis://{}",host)).map_err(|e|return Error::DBPoolError(e.to_string()));
    // let redis_client = res.unwrap().get_async_connection().await.map_err(|e|return Error::DatabaseError(e.to_string()));
    Ok(pool.unwrap())
}