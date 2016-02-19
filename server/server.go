package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {

	logFile, logErr := os.OpenFile("./server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "syn_tool start Failed")
		os.Exit(1)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	runtime.GOMAXPROCS(runtime.NumCPU())
	server := NewServer("0.0.0.0", 5779)
	server.StartServer()
}
