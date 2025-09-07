package main

import (
	"fmt"
	"levelzero/protos/stockpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type stockServer struct {
	stockpb.UnimplementedStockPublisherServer
}

// StartMarket proto method
// Makes a market each request, runs it, writes each update to the grpc stream
func (s *stockServer) StartMarket(req *stockpb.StartMarketRequest, stream grpc.ServerStreamingServer[stockpb.Stock]) error {
	startTime := time.Now()

	if req.StartDate != nil {
		fmt.Println("start date passed")
		startTime = req.StartDate.AsTime()
	}
	
	market := MakeMarket(req.Stocks, startTime)
	stockChannel := market.Run(stream.Context())

	for update := range stockChannel {
		// data changes when i write to stream. either copy everything or come up with better solution
		if err := stream.Send(update.Data); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	port := ":3131"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	stockpb.RegisterStockPublisherServer(grpcServer, &stockServer{})

	// this will let the clients inspect what services are available
	reflection.Register(grpcServer)

	fmt.Printf("Grpc service is running on %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
