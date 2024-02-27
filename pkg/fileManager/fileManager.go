package fileManager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emromerog/indexer-api/pkg/models"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
)

/*Reads all directories*/
func ReadDir() {
	log.Printf("Reading directories...")

	var wg sync.WaitGroup
	var emails []models.Email

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pathMails := "../data/maildir"

	maildirPath := filepath.Join(currentDir, pathMails)

	err = filepath.Walk(maildirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//fmt.Println(path, info.Size())
			if info.IsDir() {
				//fmt.Printf("Directorio: %s\n", path)
			} else {
				//fmt.Printf("Archivo: %s, Tamaño: %d bytes\n", path, info.Size())
				//readFileWithBufio(path)
				wg.Add(1)
				go readFile(path, &wg, &emails) //cada llamada a readFileWithBufio se ejecute en su propia goroutine y que main espere a que todas terminen antes de salir
				//wg.Wait()
			}
			return nil
		})
	if err != nil {
		//log.Println(err)
	}

	wg.Wait()

	//517425
	//slice 487005
	// Obtener la longitud del slice
	longitud := len(emails)

	// Imprimir la longitud del slice
	fmt.Printf("El slice tiene %d elementos.\n", longitud)

	zincsearch.BulkData(emails)
}

/*Reads the files found in each directory*/
func readFile(filePath string, wg *sync.WaitGroup, emails *[]models.Email) /*(models.Email, error)*/ {
	//para garantizar que la goroutine se marque como finalizada incluso si hay un error al abrir el archivo.
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error al abrir el archivo %s: %s\n", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var email models.Email

	for {
		//var email Email

		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break // Se alcanzó el final del archivo
			} else {
				log.Printf("Error al leer el archivo %s: %s\n", filePath, err)
				break
			}
		}

		switch {
		case strings.HasPrefix(line, "Message-ID:"):
			email.MessageId = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		case strings.HasPrefix(line, "Date:"):
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		case strings.HasPrefix(line, "From:"):
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		case strings.HasPrefix(line, "To:"):
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		case strings.HasPrefix(line, "Subject:"):
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		//case strings.HasPrefix(line, "Cc:"):
		//email.Cc = strings.TrimSpace(strings.TrimPrefix(line, "Cc:"))
		case strings.HasPrefix(line, "Mime-Version:"):
			email.MimeVersion = strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:"))
		case strings.HasPrefix(line, "Content-Type:"):
			email.ContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
			// Puedes agregar más casos según sea necesario
		}
	}

	*emails = append(*emails, email)
}
