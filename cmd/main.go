package main

import (
	"log"

	"github.com/emromerog/indexer-api/pkg/env"
	"github.com/emromerog/indexer-api/pkg/fileManager"
	"github.com/emromerog/indexer-api/pkg/http"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
)

func main() {
	env.LoadVars()

	existIndex, err := zincsearch.CheckIndexExists()
	checkError("Error verifying the index: ", err)

	//fmt.Println(os.Getenv("STRONGEST_AVENGER"))

	if !existIndex {
		fileManager.ReadDirectories()
	}

	err = http.InitializeServer()
	checkError("Error starting the server: ", err)
}

func checkError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
