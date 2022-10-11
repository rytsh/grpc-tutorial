package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/rytsh/grpc-tutorial/greet/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement calculator.CalculatorServer.
type server struct {
	pb.UnimplementedGreetServiceServer
}

// Sum implements calculator.CalculatorServer
func (s *server) GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	log.Printf("Received: %v", in.FirstName)
	for i := 0; i < 10; i++ {
		result := fmt.Sprintf("Hello %v", in.FirstName)

		if err := stream.Send(&pb.GreetResponse{Result: result}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	log.Printf("Received a LongGreet request")
	result := ""

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			return stream.SendAndClose(&pb.GreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		result += fmt.Sprintf("Hello %s!\n", in.FirstName)
	}
}

func (s *server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	log.Printf("Received a GreetEveryone request")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		result := fmt.Sprintf("Hello %s!\n", in.GetFirstName())
		if err := stream.Send(&pb.GreetResponse{Result: result}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
