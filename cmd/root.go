package cmd

import (
	"github.com/spf13/cobra"
	"signature-server/util"
	"sync"
)

var rootCmd = &cobra.Command{
	Version: "v1.0",
	Use:     "signature-server",
	Short:   "sign and verify digital data",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Fatalf("Unable to execute command: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

type runnable interface {
	Run()
}

func runServer(wg *sync.WaitGroup, server runnable) {
	defer wg.Done()
	server.Run()

}
