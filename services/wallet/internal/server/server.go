package server

import (
	"log"
	"net"
	"safir/libs/appconfigs"
	"safir/libs/appstates"
	"voolibow_wallet/internal/client"
	"voolibow_wallet/internal/database"
	"voolibow_wallet/internal/handlers"
	"voolibow_wallet/internal/repository"
	"voolibow_wallet/internal/services"
	pb "voolibow_wallet/proto/api"

	"google.golang.org/grpc"
)

func RunServer() {
	var (
		listenAddress = appconfigs.String("listen-address", "Server listen address")
		dbHost        = appconfigs.String("db-host", "Database host address")
		dbPort        = appconfigs.Int("db-port", "Database port")
		dbName        = appconfigs.String("db-name", "Database name")
		dbUsername    = appconfigs.String("db-username", "Database username")
		dbPassword    = appconfigs.String("db-password", "Database password")
		infuraKey     = appconfigs.String("infura-key", "Infura Key")
	)

	if err := appconfigs.Parse(); err != nil {
		appstates.PanicMissingEnvParams(err.Error())
	}
	appstates.DoneConfigExtraction("all config params received")

	db, err := database.ConnectToPostgres(*dbHost, *dbPort, *dbName, *dbUsername, *dbPassword)
	if err != nil {
		appstates.PanicDBConnectionFailed(err.Error())
	}
	appstates.DoneDbConnection("connected to database")

	infuraClient, err := client.NewETHClient(*infuraKey)
	if err != nil {
		log.Fatal(err)
	}
	var (
		repository repository.WalletRepository = repository.NewWalletRepository(db)

		walletService services.WalletService = services.NewWalletService(repository, infuraClient)

		walletHadler = handlers.NewWalletHandler(walletService)
	)

	listener, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("error : %v", err)
		appstates.PanicServerSocketFailure(err.Error())
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWalletServiceServer(grpcServer, walletHadler)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("error : %v", err)
	}

}
