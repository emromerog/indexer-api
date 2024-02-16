package http

import (
	"fmt"
	"net/http"
	//"github.com/emromerog/indexer-api/app/controllers"
)

func StartServer() {
	//http.HandleFunc("/", controllers.HelloWorldHandler)

	port := "8080"
	fmt.Printf("Iniciando el servidor en http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
