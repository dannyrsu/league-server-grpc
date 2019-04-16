package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/reflection"

	pb "github.com/dannyrsu/league-grpc-server/leagueservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Coonection uses TLS if true")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS Key file")
	port     = flag.Int("port", 50051, "The server port")
)

type leagueServer struct{}

func main() {
	// if we crash the code, we get the filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)

		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLeagueApiServer(grpcServer, &leagueServer{})
	reflection.Register(grpcServer)

	go func() {
		fmt.Println("Starting League Server ...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	grpcServer.Stop()
	lis.Close()
	fmt.Println("League Server Stopped...")
}
