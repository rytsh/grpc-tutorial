package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"strings"

	pb "github.com/rytsh/grpc-tutorial/calculator/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	sum   = flag.Bool("sum", false, "Sum use with a and b arguments")
	prime = flag.Bool("prime", false, "Prime number, use with a argument")
	avg   = flag.String("avg", "", "Average number, comma separated list of numbers")
	max   = flag.String("max", "", "Maximum number, comma separated list of numbers")
	sqrt  = flag.Int64("sqrt", 0, "square root of a number")

	a = flag.Int64("a", 0, "First number")
	b = flag.Int64("b", 0, "Second number")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	ctx := context.Background()

	if *sum {
		r, err := c.Sum(ctx, &pb.SumRequest{A: *a, B: *b})
		if err != nil {
			log.Fatalf("could not calculate: %v", err)
		}
		log.Printf("Sum: %d", r.Sum)
	}

	if *prime {
		doPrime(ctx, c, *a)
	}

	if *avg != "" {
		var numbers []int64
		for _, s := range strings.Split(*avg, ",") {
			n, _ := strconv.ParseInt(s, 10, 64)
			numbers = append(numbers, n)
		}
		log.Printf("Numbers: %v", numbers)
		doAvg(ctx, c, numbers)
	}

	if *max != "" {
		var numbers []int64
		for _, s := range strings.Split(*max, ",") {
			n, _ := strconv.ParseInt(s, 10, 64)
			numbers = append(numbers, n)
		}
		log.Printf("Numbers: %v", numbers)
		doMax(ctx, c, numbers)
	}

	if *sqrt != 0 {
		doSqrt(ctx, c, *sqrt)
	}
}
