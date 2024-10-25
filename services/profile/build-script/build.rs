fn main(){
    tonic_build::configure()
    .build_client(false)
    .out_dir("proto/api")
    .compile(&["proto/protos/profile.proto"],&["path"])
    .expect("failed to compile protos");
}