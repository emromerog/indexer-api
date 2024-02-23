package main

import (
	"fmt"
	"os"

	"github.com/emromerog/indexer-api/infrastructure/http"
)

func main() {
	err := http.InitializeServer()

	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
		os.Exit(1)
	}
}
