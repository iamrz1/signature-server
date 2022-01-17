package cmd

import (
	"github.com/spf13/cobra"
	"signature-server/api"
	"signature-server/config"
	memDB "signature-server/data/memory"
	"signature-server/util"
	"sync"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start api server",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	appCnf, err := config.AppCnf()
	if err != nil {
		util.Fatal(err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	sStore, err := memDB.NewSignatureStore(appCnf.DaemonKey)
	if err != nil {
		util.Fatal(err.Error())
		return
	}
	tStore := memDB.NewTransactionStore()
	go runServer(&wg, api.NewServer("api", appCnf.ServerPort, appCnf.Timeout, api.NewAPIRouter(sStore, tStore)))
	go runServer(&wg, api.NewServer("system", appCnf.SystemPort, appCnf.Timeout, api.NewSystemRouter()))
	wg.Wait()
}
