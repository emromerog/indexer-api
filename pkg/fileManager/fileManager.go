package fileManager

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emromerog/indexer-api/pkg/models"
	"github.com/emromerog/indexer-api/pkg/zincsearch"
)

/*Reads all directories recursively*/
func ReadDirectories() {
	log.Printf("Reading directories...")

	var wg sync.WaitGroup
	var emails []models.Email
	var mu sync.Mutex

	pathMails := "../data/maildir"

	// Canal para señalar finalización
	//done := make(chan struct{})

	err := filepath.Walk(pathMails,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				wg.Add(1)
				go readFile(path, &wg, &emails, &mu)
			}

			return nil
		})
	if err != nil {
		log.Printf("error processing directory: %v", err)
	}

	wg.Wait()

	/*go func() {
		// Esperar a que todas las goroutines terminen
		wg.Wait()

		// Cerrar el canal al finalizar
		close(done)
	}()*/

	lenEmails := len(emails)
	log.Printf("The slice has %d elements...\n", lenEmails)

	///////////////////////////////////
	//const BulkDataBatchSize = 50000

	// for i := 0; i < len(emails); i += BulkDataBatchSize {
	// 	end := i + BulkDataBatchSize
	// 	if end > len(emails) {
	// 		end = len(emails)
	// 	}

	// 	// Crear un lote de datos
	// 	batch := emails[i:end]

	// 	wg.Add(1)
	// 	go zincsearch.BulkData(batch, &wg)
	// 	/*if err != nil {
	// 		log.Printf("Error al enviar bulk data: %v\n", err)
	// 		//return fmt.Errorf("error al realizar la solicitud: %v", err)
	// 	}*/
	// }

	/////////////////////////////////////////

	wg.Wait()

	zincsearch.BulkData(emails /*, &wg*/)

	log.Println("Todas las goroutines han terminado.")
}

/*Reads the files found in each directory*/
func readFile(filePath string, wg *sync.WaitGroup, emails *[]models.Email, mu *sync.Mutex) /*(models.Email, error)*/ {
	//para garantizar que la goroutine se marque como finalizada incluso si hay un error al abrir el archivo.
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error al abrir el archivo %s: %s\n", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	//scanner := bufio.NewScanner(file)
	var email models.Email
	/*var flagContent bool*/
	//var contentBuffer strings.Builder
	var content string

	for {
		//for scanner.Scan() {
		line, err := reader.ReadString('\n')
		//line := scanner.Text()
		if err != nil {
			if err.Error() == "EOF" {
				break // Se alcanzó el final del archivo
			} else {
				log.Printf("Error al leer el archivo %s: %s\n", filePath, err)
				break
			}
		}

		switch {
		/*case strings.HasPrefix(line, "Message-ID:"):
		email.MessageId = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))*/
		case strings.HasPrefix(line, "Date:"):
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		case strings.HasPrefix(line, "From:"):
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		case strings.HasPrefix(line, "To:"):
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		case strings.HasPrefix(line, "Subject:"):
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		case strings.HasPrefix(line, "Mime-Version:"):
			email.MimeVersion = strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:"))
		case strings.HasPrefix(line, "Content-Type:"):
			email.ContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
		case strings.HasPrefix(line, "Content-Transfer-Encoding:"):
			email.ContentTransferEncoding = strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:"))
			/*case strings.HasPrefix(line, "X-FileName:"):
			//flagContent = true
			email.XFileName = strings.TrimSpace(strings.TrimPrefix(line, "X-FileName:"))
			case strings.HasPrefix(line, "Content:"):
			email.Content = strings.TrimSpace(strings.TrimPrefix(line, "Content:"))*/
		}

		/*if flagContent {
			content = content + line
		}*/

		//content = content + line
	}

	/*if err := scanner.Err(); err != nil {
		log.Printf("Error al leer el archivo %s: %s\n", filePath, err)
	}}*/

	email.Content = content

	mu.Lock()

	//email.Content = contentBuffer.String()

	*emails = append(*emails, email)

	/*defer*/
	mu.Unlock()

}

/*func sendBulkData() {

}*/
