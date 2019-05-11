package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/hazward/plexcluster/plexcluster"
	"google.golang.org/grpc"
)

type transcoderServer struct {
	jobs chan *pb.JobRequest
	l    *log.Logger
}

func (srv transcoderServer) Transcode(workerStatus *pb.WorkerStatus, stream pb.TranscoderService_TranscodeServer) error {
	done := make(chan bool)
	var err error

	go func() {
		for {
			job := <-srv.jobs
			srv.l.Printf("Sending Job %s to Worker %s\n%s", job.Id, workerStatus.WorkerId, job.Args)
			err = stream.Send(job)
			if err != nil {
				srv.l.Printf("An error occured when sending job %s to worker %s: %v", job.Id, workerStatus.WorkerId, err)
				srv.jobs <- job
				done <- true
			}
		}
	}()

	<-done
	return err
}

func (srv transcoderServer) SendJob(_ context.Context, job *pb.JobRequest) (*pb.JobStatus, error) {
	srv.l.Printf("Received job %s, adding it to queue", job.Id)
	srv.jobs <- job
	return &pb.JobStatus{JobId: job.Id, Status: pb.Status_SCHEDULED}, nil
}

func (srv transcoderServer) parseAddr(addr string) (net.Listener, error) {
	tokens := strings.Split(addr, "://")

	srv.l.Printf("Setting up server at using %s", addr)
	addressString := strings.Join(tokens[1:], "")
	switch protocol := tokens[0]; protocol {
	case "tcp":
		return net.Listen("tcp", addressString)
	case "unix":
		unixAddr, err := net.ResolveUnixAddr("unix", addressString)
		if err != nil {
			return nil, fmt.Errorf("unable to resolve unix address '%s': %v", addr, err)
		}
		return net.ListenUnix("unix", unixAddr)
	default:
		return nil, fmt.Errorf("unknown protocol: %s", addr)
	}
}

func Run(serverAddr string) {
	srv := grpc.NewServer()

	server := &transcoderServer{
		jobs: make(chan *pb.JobRequest),
		l:    log.New(os.Stdout, "[SERVER] ", log.Ldate|log.Ltime),
	}

	pb.RegisterTranscoderServiceServer(srv, server)
	l, err := server.parseAddr(serverAddr)
	if err != nil {
		server.l.Fatalf("could make listener from '%s': %v", serverAddr, err)
	}
	server.l.Printf("Transcoding server waiting on %s for jobs", serverAddr)
	server.l.Fatal(srv.Serve(l))
}
