package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// importing generated stubs
	pb "github.com/cloverway/schema/pbgo/v1/sample"

	"google.golang.org/grpc"
)

// SampleServerImpl will implement the service defined in protocol buffer definitions
type SampleServerImpl struct {
	pb.UnimplementedSampleServer
}

// SayHello is the implementation of RPC call defined in protocol definitions.
// This will take HelloRequest message and return HelloReply
func (g *SampleServerImpl) GetVersion(ctx context.Context, request *pb.GetVersionRequest) (*pb.Version, error) {
	fmt.Println("GetVersion invoked ... ")
	return &pb.Version{
		Version:     fmt.Sprintf("%s.%d.%d", "v1", 0, 1),
		CommitLong:  "This is a initial commit implement a sample grpc service",
		CommitShort: "inittial commit",
	}, nil
}

func main() {
	done := make(chan bool)
	serveGrpcSvc()
	<-done
}

func serveGrpcSvc() {
	// create new gRPC server
	server := grpc.NewServer()
	// register the SampleServerImpl on the gRPC server
	pb.RegisterSampleServer(server, &SampleServerImpl{})
	// start listening on port :8080 for a tcp connection
	if l, err := net.Listen("tcp", ":8080"); err != nil {
		log.Fatal("error in listening on port :8080", err)
	} else {
		// the gRPC server
		if err := server.Serve(l); err != nil {
			log.Fatal("unable to start server", err)
		}
	}
}

func serveRestSvc() {

}
