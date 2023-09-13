package main

import (
	"signature-server/api"
	"signature-server/config"
	"signature-server/logger"
	"sync"
)

// @title Signature Server
// @version v1.0
// @description This is a signature server
// @termsOfService tbd
// @contact.name Rezoan Tamal
// @contact.email rezoan.tamal@gmail.com
// @host localhost:8080
// @BasePath /

func main() {
	appCnf, err := config.AppCnf()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	logger.InitLogger(appCnf.LogLevel)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go runServer(&wg, api.NewServer("api", appCnf.ServerPort, appCnf.Timeout, api.NewAPIRouter(&appCnf)))
	go runServer(&wg, api.NewServer("system", appCnf.SystemPort, appCnf.Timeout, api.NewSystemRouter()))
	wg.Wait()
}

func runServer(wg *sync.WaitGroup, server *api.Server) {
	defer wg.Done()
	server.Run()
}
