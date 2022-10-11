package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"
)

func doMax(ctx context.Context, c pb.CalculatorClient, x []int64) {
	log.Println("Starting to do a Max RPC...")

	stream, err := c.Max(ctx)
	if err != nil {
		log.Fatalf("Error while calling Max RPC: %v", err)
	}

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, number := range x {
			if err := stream.Send(&pb.MaxRequest{Number: number}); err != nil {
				log.Fatalf("Error while sending request: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
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

			log.Printf("Response from Max: %v", msg.GetMax())
		}
	}()

	wg.Wait()
}
