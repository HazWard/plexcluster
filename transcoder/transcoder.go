package transcoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hazward/plexcluster/types"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

var webClient = &http.Client{
	Timeout: 10 * time.Second,
}

func handleJobRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var job types.Job
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(body, &job)
		if err != nil {
			log.Println(err)
			return
		}
		if time.Now().After(job.Expiry) {
			w.WriteHeader(403)
			_, err = fmt.Fprintf(w, "Job '%s' discard because it expired", job.ID)
			log.Println(err)
			return
		}
		_, err = fmt.Fprintf(w, "Job '%s' received successfully", job.ID)
		go runJob(job)
		log.Println(err)
	default:
		w.WriteHeader(405)
		_, err := fmt.Fprintf(w, "'%s' is not implemented, only POST is supported", r.Method)
		log.Println(err)
	}
}

func runJob(job types.Job) {
	log.Printf("Executing job '%s': %s", job.ID, job.Args)
	cmd := exec.Command("/usr/lib/plexmediaserver/plex_transcoder", job.Args...)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Job '%s' finished with error: %v", job.ID, err)
}

// Run registers a transcoder and starts up an endpoint to receive jobs from a load balancer
func Run(loadBalancerAddr string, transcoderType types.TranscoderType) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln(err)
	}

	name, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	listeningPort := listener.Addr().(*net.TCPAddr).Port

	info := types.TranscoderInfo{
		Name: name,
		Port: listeningPort,
		Type: transcoderType,
	}

	data, err := json.Marshal(info)
	if err != nil {
		log.Fatalln(err)
	}
	loadBalancerTranscoderRegistrationURL := fmt.Sprintf("http://%s/transcoders", loadBalancerAddr)
	resp, err := webClient.Post(loadBalancerTranscoderRegistrationURL, "application/json", bytes.NewReader(data))
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("could not register transcoder: %s | %v", err, resp)
	}
	transcoderKey, _ := ioutil.ReadAll(resp.Body)

	leaveHost := func() {
		loadBalancerTranscoderURLString := fmt.Sprintf("http://%s/transcoders/%s", loadBalancerAddr, string(transcoderKey))
		loadBalancerTranscoderURL, err := url.Parse(loadBalancerTranscoderURLString)
		if err != nil {
			log.Fatalf("could not generate transcoder removal URL: %s", err)
		}
		req := &http.Request{
			Method: http.MethodDelete,
			URL: loadBalancerTranscoderURL,
		}
		resp, err := webClient.Do(req)
		defer resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("could not remove transcoder: %s | %v", err, resp)
		}
	}
	defer leaveHost()
	fmt.Println("Listening on port:", listeningPort)
	http.HandleFunc("/jobs", handleJobRequest)
	log.Fatalln(http.Serve(listener, nil))
}
