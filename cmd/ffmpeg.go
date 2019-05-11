package cmd

import (
	"fmt"
	"github.com/hazward/plexcluster/ffmpeg"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var ffmpegCmd = &cobra.Command{
	Use:              "ffmpeg [TRANSCODING ARGUMENTS]",
	Short:            "Send the transcoding task to a transcoding server",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		var ok bool

		if strings.Compare(remoteServer, "") == 0 {
			remoteServer, ok = os.LookupEnv("REMOTE_TRANSCODER_SERVER")
			if !ok {
				fmt.Println("--server cannot be blank")
				fmt.Println(cmd.UsageString())
				os.Exit(2)
			}
		}

		env := os.Environ()
		ffmpeg.Run(remoteServer, args, env)
	},
}

func init() {
	ffmpegCmd.PersistentFlags().StringVarP(&remoteServer, "server", "", "", "Server address of transcoding server")
}
