package server

import (
	"context"
	"fmt"
	"net"

	"buf.build/go/protovalidate"
	"github.com/Auction-Application/be-user-service/internal/database"
	"github.com/Auction-Application/be-user-service/internal/database/usersTableQuery"
	userPb "github.com/Auction-Application/be-user-service/rpc/gen/user/v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type config struct {
	port int
}

type UserServer struct {
	userPb.UnimplementedUserServiceServer
	databaseConn *pgx.Conn
}

func (s *UserServer) InsertUser(ctx context.Context, req *userPb.InsertUserPayload) (*userPb.InsertUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user credentials")
	}

	usersTable := usersTableQuery.New(s.databaseConn)
	err := usersTable.InsertUser(ctx, usersTableQuery.InsertUserParams{Username: *req.User.Username, Email: *req.User.Email, AuthID: *req.User.AuthID})
	if err != nil {
		// return &userPb.InsertUserResponse{Success: proto.Bool(false), Message: proto.String(err.Error())}, nil
		return nil, status.Error(codes.Internal, "error inserting user")
	}
	return &userPb.InsertUserResponse{Success: new(true), Message: new("User Inserted")}, nil

}

func NewUserServer() *UserServer {
	return &UserServer{
		databaseConn: database.ConnectToDB(),
	}
}

func Run(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Println("failed to start the server")
		return err
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	ctx := context.Background()
	userServer := NewUserServer()

	defer userServer.databaseConn.Close(ctx)

	userPb.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Println("error listening the server")
		return err
	}
	return nil
}
