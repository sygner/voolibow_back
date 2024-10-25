// Package client provides functionalities for establishing a gRPC client connection.
package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcServerConnection establishes a connection to a gRPC server.
// It takes the server_address as a parameter and returns a *grpc.ClientConn.
// Insecure credentials are used for simplicity in this example, which is suitable for development or testing environments.
func GrpcServerConnection(server_address string) *grpc.ClientConn {
	// Dial connects to the gRPC server using the specified server_address and insecure credentials.
	conn, err := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// If an error occurs during connection, log the error and terminate the program.
		log.Fatalf("did not connect: %v", err)
	}
	// Return the established gRPC client connection.
	return conn
}
