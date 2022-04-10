package main

import (
	"context"
	"log"
	"time"

	pb "github.com/cloverway/schema/pbgo/v1/echoserver"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewEchoServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.GetEchoMsg(ctx, &pb.EchoRequest{
		Msg: "Client echo message",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Res: %s", r)
}
