use sqlx::{postgres::PgPoolOptions, Postgres};
use super::types::error::Error;
pub async fn postgres_connection(username:String,password:String,host:String,port:usize,table:String)->Result<sqlx::Pool<Postgres>,Error>{
    let uri = format!("postgres://{}:{}@{}:{}/{}",username,password,host,port,table);
    let pool = PgPoolOptions::new()
    .max_connections(5)
    .connect(uri.as_str()).await.map_err(|e|return Error::DBPoolError(e.to_string()));
    Ok(pool.unwrap())
}
