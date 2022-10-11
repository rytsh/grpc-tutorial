package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/rytsh/grpc-tutorial/greet/proto"
)

func doGreetEveryone(ctx context.Context, c pb.GreetServiceClient) {
	log.Println("Starting to do a GreetEveryone RPC...")

	reqs := []*pb.GreetRequest{
		{FirstName: "Selin"},
		{FirstName: "Melda"},
	}

	stream, err := c.GreetEveryone(ctx)
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, req := range reqs {
			log.Printf("Sending req: %v", req)
			if err := stream.Send(req); err != nil {
				log.Fatalf("Error while sending request: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// we've reached the end of the stream
				break
			}
			if err != nil {
				log.Fatalf("Error while reading stream: %v", err)
			}

			log.Printf("Response from GreetEveryone: %v", msg.GetResult())
		}
	}()

	wg.Wait()
}
