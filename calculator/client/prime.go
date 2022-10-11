package main

import (
	"context"
	"io"
	"log"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"
)

func doPrime(ctx context.Context, c pb.CalculatorClient, x int64) {
	log.Println("Starting to do a Prime RPC...")

	req := &pb.PrimeRequest{
		Number: x,
	}

	resStream, err := c.Prime(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling Prime RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from Prime: %v", msg.GetDivide())
	}

}
