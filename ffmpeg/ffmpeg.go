package ffmpeg

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	pb "github.com/hazward/plexcluster/plexcluster"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type transcoderServer struct {

}

func (srv transcoderServer) Transcode(pb.TranscoderService_TranscodeServer) error {
	panic("implement me")
}

func runJob(job pb.JobRequest) {
	log.Printf("Executing job '%s': %s", job.Id, job.Args)
	cmd := exec.Command("/usr/lib/plexmediaserver/plex_transcoder", job.Args...)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	if err != nil {
		log.Printf("Job '%s' finished with error: %v", job.Id, err)
		return
	}
	log.Printf("Job '%s' finished successfully", job.Id)
}

func Run(serverAddr string, args []string) {
	h := sha256.New()
	h.Write([]byte(strings.Join(args, "")))
	job := pb.JobRequest{
		Id: fmt.Sprintf("%s", string(h.Sum(nil)[:8])),
		Args: args,
		Expiry: time.Now().Add(1 * time.Minute).Unix(),
	}

	srv := grpc.NewServer()
	var transcodes transcoderServer
	tokens := strings.Split(serverAddr, "://")
	pb.RegisterTranscoderServiceServer(srv, transcodes)
	l, err := net.Listen(tokens[0], tokens[1])
	if err != nil {
		log.Printf("could not listen to %s: %v", tokens[1], err)
		log.Print("using local transcoder instead")
		runJob(job)
	} else {
		err = srv.Serve(l)
		if err != nil {
			runJob(job)
		}
	}

	time.Sleep(1 * time.Hour)
}
