package common

import (
	"fmt"
	pb "github.com/hazward/plexcluster/plexcluster"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const TranscoderPath = "/usr/lib/plexmediaserver/plex_transcoder"
const defaultPlexHost = "127.0.0.1:32400"

var simpleArgumentRegex = regexp.MustCompile(`(?m)[A-Za-z0-9]+`)

func RunJob(l *log.Logger, job pb.JobRequest) {
	transcoderPath, ok := os.LookupEnv("TRANSCODER_PATH")
	if !ok {
		transcoderPath = TranscoderPath
	}

	args := flattenArguments(replacePlexHostArguments(job.Args))

	l.Printf("Executing job '%s': %s", job.Id, args)
	command := exec.Command(transcoderPath, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Env = job.Env
	fmt.Println("Environment:", job.Env)
	l.Printf("Running command and waiting for it to finish...")
	err := command.Start()
	if err != nil {
		l.Printf("Job '%s' unable to start: %v", job.Id, err)
		return
	}
	err = command.Wait()
	if err != nil {
		l.Printf("Job '%s' finished with error: %v", job.Id, err)
		return
	}
	l.Printf("Job '%s' finished successfully", job.Id)
}

func replacePlexHostArguments(args []string) []string {
	for i, arg := range args {
		if strings.Contains(arg, defaultPlexHost) {
			host, ok := os.LookupEnv("PLEX_HOST")
			if ok {
				args[i] = strings.ReplaceAll(arg, defaultPlexHost, host)
			}
		}
	}
	return args
}

func flattenArguments(args []string) []string {
	var finalArgs []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			finalArgs = append(finalArgs, arg)
			newIndex, combinedArgs := combineArgs(i+1, args)
			if strings.Compare(combinedArgs, "") != 0 {
				if strings.Contains(arg, "loglevel") {
					finalArgs = append(finalArgs, "debug")
				} else {
					finalArgs = append(finalArgs, sanitizeArguments(combinedArgs))
				}
			}
			i = newIndex
			continue
		}
		finalArgs = append(finalArgs, fmt.Sprintf(`%s`, arg))
		i++
	}
	return finalArgs
}

func combineArgs(i int, args []string) (int, string) {
	if i >= len(args) {
		return len(args), ""
	}

	if strings.HasPrefix(args[i], "-") {
		return i + 1, args[i]
	}
	var combined strings.Builder

	for i < len(args) && !strings.HasPrefix(args[i], "-") {
		combined.WriteString(fmt.Sprintf("%s ", args[i]))
		i++
	}

	return i, strings.TrimSpace(combined.String())
}

// sanitizeArguments is to add quotes
// around values that need it
func sanitizeArguments(arg string) string {
	matches := simpleArgumentRegex.FindAllString(arg, -1)
	if strings.HasPrefix(arg, "-") {
		return arg
	}

	if len(matches) > 1 {
		return fmt.Sprintf(`"%s"`, arg)
	}

	return arg
}