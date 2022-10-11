package main

import (
	"context"
	"log"
	"time"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"
)

func doAvg(ctx context.Context, c pb.CalculatorClient, x []int64) {
	log.Println("Starting to do a Avg RPC...")

	stream, err := c.Avg(ctx)
	if err != nil {
		log.Fatalf("Error while calling Avg RPC: %v", err)
	}

	for _, number := range x {
		if err := stream.Send(&pb.AvgRequest{Number: number}); err != nil {
			log.Fatalf("Error while sending request: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	log.Printf("Response from Avg: %v", res.GetNumber())
}
