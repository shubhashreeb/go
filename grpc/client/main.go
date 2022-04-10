package main

import (
	"context"
	"log"
	"time"

	pb "github.com/cloverway/schema/pbgo/v1/sample"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewSampleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetVersion(ctx, &pb.GetVersionRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("URL: %s", r)
}
