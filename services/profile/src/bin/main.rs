use std::sync::Arc;
use profilelib::app::config::ParseConfig;
use profilelib::app::{database, handlers};
use profilelib::profile_proto::profile_service_server::ProfileServiceServer;
use clap::Parser;
use tonic::transport::Server;

#[tokio::main]
async fn main() {
    
    println!("[INFO] Parse Input Config");
    let parsed =  ParseConfig::parse();
    println!("[INFO] Connecting To PostgresDB...");
    let pg_db_pool = database::postgres_connection(parsed.db_username, parsed.db_password, parsed.db_host, parsed.db_port,parsed.db_name)
    .await.unwrap();
    println!("[INFO] Connected To PostgresDB!");

    let pg_db_pool = Arc::new(pg_db_pool);
    
    println!("[INFO] Init Services");
    // init services


    let profile_handler = handlers::profile::ProfileHandler::new(pg_db_pool);
    let profile_service= ProfileServiceServer::new(profile_handler);  

    println!("[INFO] Running Server On {}",parsed.listen_address);
    Server::builder()
    .add_service(profile_service)
    .serve(parsed.listen_address.parse().expect("could not parse the listener address"))
    .await
    .unwrap()
}
