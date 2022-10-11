package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement calculator.CalculatorServer.
type server struct {
	pb.UnimplementedCalculatorServer
}

// Sum implements calculator.CalculatorServer
func (s *server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Received: A: %v, B: %v", in.A, in.B)
	return &pb.SumResponse{Sum: in.A + in.B}, nil
}

func (s *server) Prime(in *pb.PrimeRequest, stream pb.Calculator_PrimeServer) error {
	N := in.Number
	log.Printf("Received: %v", N)

	k := int64(2)
	for N > 1 {
		if N%k == 0 { // if k evenly divides into N
			N /= k // divide N by k so that we have the rest of the number left.
			stream.Send(&pb.PrimeResponse{Divide: k})
		} else {
			k++
		}
	}

	return nil
}

func (s *server) Avg(stream pb.Calculator_AvgServer) error {
	var sum int64
	var count int64
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			return stream.SendAndClose(&pb.AvgResponse{Number: float64(sum) / float64(count)})
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Received: %v", in.GetNumber())
		sum += in.GetNumber()
		count++
	}
}

func (s *server) Max(stream pb.Calculator_MaxServer) error {
	var max int64
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Received: %v", in.GetNumber())
		if in.GetNumber() > max {
			max = in.GetNumber()
			if err := stream.Send(&pb.MaxResponse{Max: max}); err != nil {
				return err
			}
		}
	}
}

func (s *server) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	n := in.GetNumber()
	log.Printf("Received: %v", in.GetNumber())

	if n < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", n),
		)
	}

	return &pb.SqrtResponse{Number: math.Sqrt(float64(n))}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
