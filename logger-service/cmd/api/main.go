package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/cmd/api/config"
	"logger-service/cmd/api/grpctest"
	"logger-service/cmd/api/router"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {

	var config = config.ConfigParam{
		DBURL:      "mongodb://dockompose-mongodb-1:27017",
		DBUsername: "admin",
		DBPassword: "admin",
		WebPort:    80,
		RPCPort:    5001,
		GRPCPort:   50001,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := config.CreateConnection(ctx)
	if err != nil {
		log.Panic("Got error creating mongo db connection")
	}
	if client == nil {
		log.Panic("client is nil")
	}
	log.Println("Connected to MongoDB")

	defer func() {

		if err = client.Disconnect(ctx); err != nil {
			fmt.Println("error disconnecting")
		} else {
			fmt.Println("successfully disconnected")
		}
	}()
	err = rpc.Register(new(router.RPCServer))
	go rpcListen()
	go grpctest.GrpcListen()
	r := router.GetRouter(client)
	http.ListenAndServe(":80", r)

}

func rpcListen() error {
	fmt.Println("starting rpc server on port: ", 5001)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 5001))
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}
