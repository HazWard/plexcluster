package main

import (
	"flag"
	"fmt"
	"github.com/hazward/plexcluster/ffmpeg"
	"github.com/hazward/plexcluster/transcoder"
	pb "github.com/hazward/plexcluster/plexcluster"
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
	ffmpegServerURL := ffmpegCommand.String("server", "", "Address for GRPC server")
	transcoderCommand := flag.NewFlagSet("transcode", flag.ExitOnError)
	transcoderServerURL := transcoderCommand.String("server", "", "Address for GRPC server")

	if len(os.Args) < 3 {
		printUsage(os.Args[0], ffmpegCommand, transcoderCommand)
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
	default:
		log.Printf("%q is not valid command.\n", os.Args[1])
		flag.Usage()
		os.Exit(2)
	}

	if ffmpegCommand.Parsed() {
		ffmpeg.Run(*ffmpegServerURL, ffmpegCommand.Args())
		return
	}
	if transcoderCommand.Parsed() {
		claim := strings.TrimSpace(os.Getenv("PLEX_CLAIM"))
		transcoderType := pb.MachineType_BAREMETAL
		if strings.Compare(claim, "") != 0 {
			transcoderType = pb.MachineType_BAREMETAL
		}
		transcoder.Run(*transcoderServerURL, transcoderType)
		return
	}
}
