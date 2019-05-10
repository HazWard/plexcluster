package ffmpeg

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/hazward/plexcluster/types"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func runJob(job types.Job) {
	log.Printf("Executing job '%s': %s", job.ID, job.Args)
	cmd := exec.Command("/usr/lib/plexmediaserver/plex_transcoder", job.Args...)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	if err != nil {
		log.Printf("Job '%s' finished with error: %v", job.ID, err)
		return
	}
	log.Printf("Job '%s' finished successfully", job.ID)
}

// Run submits the current transcoding arguments args to the load balancer
// for it to schedule the job. If the load balancer can't process the job,
// the transcoding task is performed locally
func Run(queueAddr string, args []string) {
	h := sha256.New()
	h.Write([]byte(strings.Join(args, "")))
	job := types.Job{
		ID: fmt.Sprintf("%s", string(h.Sum(nil)[:8])),
		Args: args,
		Expiry: time.Now().Add(1 * time.Minute),
	}

	data, err := json.Marshal(job)
	if err != nil {
		log.Fatalln(err)
	}
	loadBalancerTranscoderURL := fmt.Sprintf("http://%s/jobs", loadBalancerAddr)
	resp, err := webClient.Post(loadBalancerTranscoderURL, "application/json", bytes.NewReader(data))
	if err != nil || resp.StatusCode != http.StatusOK{
		log.Fatalf("could submit job to loadbalancer, using local transcoder instead: %s | %v", err, resp)
		runJob(job)
	}
	time.Sleep(1 * time.Hour)
}
