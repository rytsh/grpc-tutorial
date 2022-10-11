package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/rytsh/grpc-tutorial/greet/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// doGreetManyTimes(ctx, c)
	// doLongGreet(context.Background(), c)
	doGreetEveryone(context.Background(), c)
}
