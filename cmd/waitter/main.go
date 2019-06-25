package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/john-medeiros/waitter"
	"github.com/natefinch/lumberjack"
)

var globalConfig, _ = LoadConfiguration("../../config/waitter.json")

//
func init() {
	file, err := waitter.LoadConfiguration("../../config/waitter.json")
	if err != nil {
		fmt.Println("Erro ao carregar configurações.")
		os.Exit(1)
	}

	globalConfig = file

	log.SetOutput(&lumberjack.Logger{
		Filename:   globalConfig.Logging.Path,
		MaxSize:    globalConfig.Logging.MaxSize, // megabytes
		MaxBackups: globalConfig.Logging.MaxBackups,
		MaxAge:     globalConfig.Logging.MaxAge,   //days
		Compress:   globalConfig.Logging.Compress, // disabled by default
	})

}

// QueueExec armazena todas as execuções
type QueueExec struct {
	RunFileConfig FileConfig
	RunFileList   FileList
}

func main() {

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

}
