package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	// importing generated stubs
	pb "github.com/cloverway/schema/pbgo/v1/sample"
	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	go serveGrpcSvc()
	go serveRestSvc()
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

func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in the allowed header, don't send the header
	return s, false
}

func serveRestSvc() { // creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// convert header in response(going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}))
	// setting up a dail up for gRPC service by specifying endpoint/target url
	// err := gen.RegisterGreeterHandlerFromEndpoint(context.Background(), mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})

	err := pb.RegisterSampleHandlerFromEndpoint(context.Background(), mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal(err)
	}
	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(mux),
	}
	// creating a listener for server
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	// start server
	err = server.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}
