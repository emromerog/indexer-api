package main

import (
	"log"
	"os"

	"github.com/emromerog/indexer-api/pkg/fileManager"
	"github.com/emromerog/indexer-api/pkg/http"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
)

func main() {

	existIndex, err := zincsearch.CheckIndexExists()
	checkError("Error verifying the index: ", err)

	if !existIndex {
		fileManager.ReadDir()
	}

	err = http.InitializeServer()
	checkError("Error starting the server: ", err)
}

func checkError(msg string, err error) {
	if err != nil {
		log.Printf(msg, err)
		os.Exit(1)
	}
}
