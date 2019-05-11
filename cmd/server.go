package cmd

import (
	"fmt"
	"github.com/hazward/plexcluster/server"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var serverCmd = &cobra.Command{
	Use:              "server",
	Short:            "Start a transcoding server",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		var ok bool
		if strings.Compare(remoteServer, "") == 0 {
			serverURI, ok = os.LookupEnv("SERVER_URI")
			if !ok {
				fmt.Println("--server cannot be blank")
				fmt.Println(cmd.UsageString())
				os.Exit(2)
			}
		}
		server.Run(serverURI)
	},
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverURI, "host", "", "", "Host URI to listen for workers and job submissions")
}

