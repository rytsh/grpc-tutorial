package main

import (
	"context"
	"log"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doSqrt(ctx context.Context, c pb.CalculatorClient, x int64) {
	res, err := c.Sqrt(ctx, &pb.SqrtRequest{Number: x})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Printf("RPC failed: %v", s.Message())
			log.Printf("Error code: %v", s.Code())

			if s.Code() == codes.InvalidArgument {
				log.Printf("Did you pass a negative number?")
			}

			return
		}

		log.Fatalf("Error while calling Sqrt RPC: %v", err)
		return
	}

	log.Printf("Response from Sqrt: %v", res.GetNumber())
}
