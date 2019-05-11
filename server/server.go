package server

import (
	"fmt"
	pb "github.com/hazward/plexcluster/plexcluster"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type transcoderServer struct {

}

func (srv transcoderServer) Transcode(stream pb.TranscoderService_TranscodeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func parseAddr(addr string) (net.Listener, error) {
	tokens := strings.Split(addr, "://")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid format for address: %s", addr)
	}

	if strings.Compare(tokens[0], "tcp") != 0 {
		return nil, fmt.Errorf("non tcp address provided: %s", addr)
	}
	log.Printf("Setting up server at %s", tokens[1])
	return net.Listen("tcp", tokens[1])
}

func Run(serverAddr string) {
	srv := grpc.NewServer()
	pb.RegisterTranscoderServiceServer(srv, &transcoderServer{})
	l, err := parseAddr(serverAddr)
	if err != nil {
		log.Fatalf("could make listener from '%s': %v", serverAddr, err)
	}
	log.Fatal(srv.Serve(l))
}
