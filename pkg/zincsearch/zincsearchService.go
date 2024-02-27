package zincsearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/emromerog/indexer-api/pkg/models"
	"github.com/emromerog/indexer-api/pkg/utils"
)

const (
	baseApiUrl    = "http://localhost:4080/api/"
	bulkv2Url     = "_bulkv2"
	searchUrl     = "/_search"
	existIndexUrl = "index/"
	userAuth      = "admin"
	passwordAuth  = "Complexpass#123"
)

func setBasicAuth(req *http.Request) {
	// Codificar las credenciales en base64
	auth := base64.StdEncoding.EncodeToString([]byte(userAuth + ":" + passwordAuth))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Add("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")
	req.Header.Add("Authorization", "Basic "+auth)
}

func BulkData(records []models.Email) error {
	log.Println("Sending data...")

	bulkApiURL := baseApiUrl + bulkv2Url

	body := BulkDataRequest{
		IndexName: utils.IndexName,
		Records:   records,
	}

	jsonData, err := utils.ConvertToJson(body)
	if err != nil {
		return fmt.Errorf("error al convertir a JSON: %v", err)
	}

	// Crear una solicitud POST
	req, err := http.NewRequest("POST", bulkApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud: %v", err)
	}

	setBasicAuth(req)

	// Cliente HTTP personalizado (puedes ajustar el timeout según tus necesidades)
	client := &http.Client{}

	// Realizar la solicitud POST
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta de la API no exitosa. Código de estado: %d", resp.StatusCode)
	}
	/*responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer la respuesta de la API: %v", err)
	}*/

	log.Println("Data uploaded...")

	return nil
}

func SearchData(term string, searchType string) ([]models.Email, error) {
	log.Println("Looking for information...")

	/*now := time.Now()
	startTime := now.AddDate(0, 0, -7).Format("2024-01-02T15:04:05Z")
	endTime := now.Format("2000-01-02T15:04:05Z")*/

	searchApiURL := baseApiUrl + utils.IndexName + searchUrl

	body := SearchDataRequest{
		SearchType: searchType,
		Query: SearchQuery{
			Term: term,
			//Field: "_all",
			StartTime: "2000-06-02T14:28:31.894Z",
			EndTime:   "2025-12-02T15:28:31.894Z",
		},
		From:       0,
		MaxResults: 1000,
		//Source:     []
	}

	jsonData, err := utils.ConvertToJson(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", searchApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud: %v", err)
	}

	setBasicAuth(req)

	client := &http.Client{}

	// Realizar la solicitud POST
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la solicitud: %v", err)
	}

	defer resp.Body.Close()

	// Leer la respuesta de la API
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta de la API: %v", err)
	}

	var searchDataResponse SearchDataResponse
	if err := json.Unmarshal(responseBody, &searchDataResponse); err != nil {
		return nil, fmt.Errorf("error al analizar la respuesta JSON: %v", err)
	}

	emails := make([]models.Email, 0)

	for _, hit := range searchDataResponse.Hits.Hits {
		emails = append(emails, hit.Source)
	}

	return emails, nil
}

func CheckIndexExists() (bool, error) {
	indexApiURL := baseApiUrl + existIndexUrl + utils.IndexName

	req, err := http.NewRequest("GET", indexApiURL, nil)
	if err != nil {
		return false, fmt.Errorf("error al crear la solicitud: %v", err)
	}

	setBasicAuth(req)

	// Realizar la solicitud GET
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error al realizar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error al leer la respuesta de la API: %v", err)
	}

	// Verificar el código de estado de la respuesta
	if resp.StatusCode == http.StatusOK {
		log.Println("Index with name " + utils.IndexName + " exists...")
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		log.Println("Index with name " + utils.IndexName + " does not exist...")
		return false, nil
	}
	return false, fmt.Errorf("código de estado inesperado: %d", resp.StatusCode)
}

/*func createSearchDataRequest() {

}*/
