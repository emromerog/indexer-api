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
	//pathMails := "D:/Erika/DECARGAS/enron_mail_20110402/maildir"

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
		log.Printf("Error processing directory: %v", err)
	}

	wg.Wait()

	lenEmails := len(emails)
	log.Printf("The slice has %d elements...\n", lenEmails)

	sendData(emails)
}

/*Reads the files found in each directory*/
func readFile(filePath string, wg *sync.WaitGroup, emails *[]models.Email, mu *sync.Mutex) {
	//To ensure that the goroutine is marked as completed even if there is an error opening the file.
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening the file %s: %s\n", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var email models.Email
	var flagContent bool
	var content string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				//The end of the file has been reached
				break
			} else {
				log.Printf("Error reading file %s: %s\n", filePath, err)
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
		/*case strings.HasPrefix(line, "Mime-Version:"):
			email.MimeVersion = strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:"))
		case strings.HasPrefix(line, "Content-Type:"):
			email.ContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
		case strings.HasPrefix(line, "Content-Transfer-Encoding:"):
			email.ContentTransferEncoding = strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:"))*/
		case strings.HasPrefix(line, "X-FileName:"):
			flagContent = true
			email.XFileName = strings.TrimSpace(strings.TrimPrefix(line, "X-FileName:"))
		}

		if flagContent {
			content = content + line
		}
	}

	email.Content = content

	//Ensure that only one goroutine at a time has access to a critical section of code.
	mu.Lock()

	*emails = append(*emails, email)

	//When the goroutine has completed its work on the critical section, it calls mu.Unlock() to release the mutex lock.
	mu.Unlock()
}

/*Send records in batches*/
func sendData(emails []models.Email) {
	var wg sync.WaitGroup

	batchSize := 50000

	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}

		batch := emails[i:end]

		wg.Add(1)

		go zincsearch.BulkData(batch, &wg)
		wg.Wait()
	}
}
