package transcoder

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hazward/plexcluster/common"
	pb "github.com/hazward/plexcluster/plexcluster"
	"google.golang.org/grpc"
)

func transcode(ctx context.Context, logger *log.Logger, client pb.TranscoderServiceClient) error {
	logger.Println("Connected to transcoding server, awaiting jobs...")
	id, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to setup transcoder: %v", err)
	}

	stream, err := client.Transcode(ctx, &pb.WorkerStatus{WorkerId: id})
	if err != nil {
		return fmt.Errorf("error while job stream: %v", err)
	}

	for {
		job, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return fmt.Errorf("error while getting job: %v", err)
		}
		common.RunJob(logger, *job)
	}
	return nil
}

// Run registers a transcoder and waits to receive jobs from job queue
func Run(transcodeServer string) {
	logger := log.New(os.Stdout, "[TRANSCODER] ", log.Ldate|log.Ltime)

	connection, err := grpc.Dial(transcodeServer, grpc.WithInsecure())
	if err != nil {
		logger.Fatalln(err)
	}

	ctx := context.Background()
	client := pb.NewTranscoderServiceClient(connection)
	log.Fatal(transcode(ctx, logger, client))
}
