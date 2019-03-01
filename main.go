package main

import (
	"flag"
	"fmt"
	"github.com/hazward/plexcluster/ffmpeg"
	"github.com/hazward/plexcluster/loadbalancer"
	"github.com/hazward/plexcluster/transcoder"
	"github.com/hazward/plexcluster/types"
	"log"
	"os"
	"strings"
)

func printUsage(programName string, flagsSet... *flag.FlagSet) {
	fmt.Printf("Usage: %s [command] [flags] [transcode arguments]\n\n", programName)
	fmt.Println("Commands:")
	for _, command := range flagsSet {
		fmt.Printf("%s\n", command.Name())
		command.PrintDefaults()
		fmt.Println()
	}
}

func main() {
	ffmpegCommand := flag.NewFlagSet("ffmpeg", flag.ExitOnError)
	ffmpegLbHost := ffmpegCommand.String("loadbalancer", "localhost:4545", "host:port for transcoding load balancer")
	transcoderCommand := flag.NewFlagSet("transcode", flag.ExitOnError)
	transcoderLbHost := transcoderCommand.String("loadbalancer", "localhost:4545", "host:port for transcoding load balancer")
	loadbalancerCommand := flag.NewFlagSet("loadbalancer", flag.ExitOnError)
	lbPort := loadbalancerCommand.Int("port", 4545, "Port of load balancer")

	if len(os.Args) < 3 {
		printUsage(os.Args[0], ffmpegCommand, transcoderCommand, loadbalancerCommand)
		os.Exit(2)
	}

	subCommand := os.Args[1]
	switch subCommand {
	case "ffmpeg":
		err := ffmpegCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("ffmpeg with:", os.Args[2:])
	case "transcode":
		err := transcoderCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("transcode with:", os.Args[2:])
	case "loadbalancer":
		err := loadbalancerCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("loadbalancer on %d with: %s", *lbPort, os.Args[2:])
	default:
		log.Printf("%q is not valid command.\n", os.Args[1])
		flag.Usage()
		os.Exit(2)
	}

	if ffmpegCommand.Parsed() {
		ffmpeg.Run(*ffmpegLbHost, ffmpegCommand.Args())
		return
	}
	if transcoderCommand.Parsed() {
		claim := strings.TrimSpace(os.Getenv("PLEX_CLAIM"))
		transcoderType := types.BareMetal
		if strings.Compare(claim, "") != 0 {
			transcoderType = types.Docker
		}
		transcoder.Run(*transcoderLbHost, transcoderType)
		return
	}

	if !loadbalancerCommand.Parsed() {
		err := loadbalancerCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("loadbalancer on %d with: %s", *lbPort, os.Args[2:])
	}
	loadbalancer.Run(*lbPort)
}
