package main

import (
	"fmt"
	"net"

	"github.com/Auction-Application/be-user-service/internal/server"
	userPb "github.com/Auction-Application/be-user-service/rpc/gen/user/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
		fmt.Println(err)
		return
	}

	fmt.Println("gRPC server up and ready...")
	err := Run(9001, server.NewUserServer())
	if err != nil {
		fmt.Println("Cannot run the server")
		return
	}

}

func Run(port int, userServer *server.UserServer) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Println("failed to start the server")
		return err
	}

	grpcServer := grpc.NewServer()

	userPb.RegisterUserServiceServer(grpcServer, userServer)

	fmt.Println("Listening...")

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Println("error listening the server")
		return err
	}

	fmt.Println("Hello")
	return nil
}
