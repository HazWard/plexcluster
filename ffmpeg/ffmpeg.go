package ffmpeg

import (
	"context"
	"crypto/sha256"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hazward/plexcluster/common"
	pb "github.com/hazward/plexcluster/plexcluster"
)

func isErrorStatus(status *pb.JobStatus) bool {
	return status != nil && status.Status == pb.Status_ERROR
}

func sendJob(ctx context.Context, logger *log.Logger, client pb.TranscoderServiceClient, job pb.JobRequest) error {
	logger.Printf("Sending job to server: %v", job)
	status, err := client.SendJob(ctx, &job)
	if err != nil || isErrorStatus(status) {
		return fmt.Errorf("unable to setup transcoder: %v", err)
	}
	return nil
}

func Run(serverAddr string, args, env []string) {
	logger := log.New(os.Stdout, "[FFMPEG] ", log.Ldate|log.Ltime)

	h := sha256.New()
	h.Write([]byte(strings.Join(args, "")))
	job := pb.JobRequest{
		Id:     fmt.Sprintf("%x", h.Sum(nil)[:8]),
		Args:   args,
		Expiry: time.Now().Add(1 * time.Minute).Unix(),
		Env:    env,
	}

	ctx := context.Background()
	connection, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatalln(err)
	}
	client := pb.NewTranscoderServiceClient(connection)
	err = sendJob(ctx, logger, client, job)
	if err != nil {
		logger.Printf("error while sending job, running locally instead: %s", err)
		common.RunJob(logger, job)
	}
	time.Sleep(1 * time.Hour)
}
