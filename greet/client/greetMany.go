package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/rytsh/grpc-tutorial/greet/proto"
)

func doGreetManyTimes(ctx context.Context, c pb.GreetServiceClient) {
	log.Println("Starting to do a GreetManyTimes RPC...")

	req := &pb.GreetRequest{
		FirstName: "Selin",
	}

	resStream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
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

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doLongGreet(ctx context.Context, c pb.GreetServiceClient) {
	log.Println("Starting to do a LongGreet RPC...")

	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v", err)
	}

	requests := []*pb.GreetRequest{
		{FirstName: "Selin"},
		{FirstName: "John"},
		{FirstName: "Lucy"},
		{FirstName: "Mark"},
	}

	for _, req := range requests {
		log.Printf("Sending req: %v", req)
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending request: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	log.Printf("LongGreet Response: %v", res.GetResult())
}
