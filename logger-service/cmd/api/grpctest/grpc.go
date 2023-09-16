package grpctest

import (
	"context"
	"errors"
	"fmt"
	"logger-service/cmd/api/config"
	"logger-service/cmd/api/data"
	"logger-service/cmd/api/grpcproto"

	"net"

	"google.golang.org/grpc"
)

type LogService struct {
	grpcproto.UnimplementedLogServiceServer
	Models data.LogEntry
}

func (l *LogService) WriteLog(ctx context.Context, req *grpcproto.LogRequest) (*grpcproto.LogResponse, error) {
	input := req.GetLogEntry()
	fmt.Println("go input" + input.Name + " " + input.Data)

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	dbclient, _ := config.GetClient()
	res := l.Models.Insert(logEntry, dbclient)
	if res == "" {
		err := errors.New("got error")
		return &grpcproto.LogResponse{Result: err.Error()}, err
	} else {
		return &grpcproto.LogResponse{Result: res}, nil
	}

}

func GrpcListen() {
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		fmt.Println("not able to start listener")
	}
	s := grpc.NewServer()
	grpcproto.RegisterLogServiceServer(s, &LogService{Models: data.LogEntry{Name: "ashu", Data: "ashu"}})
	fmt.Println("grpc service started on port 5050")

	if err := s.Serve(listener); err != nil {
		fmt.Println("failed to start service")
	}
}
