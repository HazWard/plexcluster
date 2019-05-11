package cmd

import "github.com/spf13/cobra"

// Global variables
var remoteServer string
var serverURI string



var rootCmd = &cobra.Command{
	Use:   "plexcluster",
	Short: "plexcluster is a drop-in replacement for Plex's transcoder with support for remote transcoding",
}

func init() {
	rootCmd.AddCommand(ffmpegCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(transcoderCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
