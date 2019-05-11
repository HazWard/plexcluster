package cmd

import (
	"fmt"
	"github.com/hazward/plexcluster/transcoder"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var transcoderCmd = &cobra.Command{
	Use:              "transcoder",
	Short:            "Subscribe the current machine as a worker for transcoding server",
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
		transcoder.Run(remoteServer)
	},
}

func init() {
	transcoderCmd.PersistentFlags().StringVarP(&remoteServer, "server", "", "", "Server address of transcoding server")
}

