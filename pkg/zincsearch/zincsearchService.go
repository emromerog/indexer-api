package zincsearch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/emromerog/indexer-api/pkg/models"
	"github.com/emromerog/indexer-api/pkg/utils"
)

var (
	//baseApiUrl    = "http://localhost:4080/api/"
	baseApiUrl    = os.Getenv("ZINCSEARCH_API_URL")
	bulkv2Url     = "_bulkv2"
	searchUrl     = "/_search"
	existIndexUrl = "index/"
	//userAuth      = os.Getenv("STRONGEST_AVENGER")
	//passwordAuth  = "Complexpass#123"
)

/*Add basic HTTP authentication headers to an HTTP request*/
func setBasicAuth(req *http.Request) {
	//Encode credentials in base64
	auth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("ZINC_FIRST_ADMIN_USER") + ":" + os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Add("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")
	req.Header.Add("Authorization", "Basic "+auth)
}

/*Submit bulk data to zincsearch for indexing*/
func BulkData(records []models.Email, wg *sync.WaitGroup) error {
	log.Println("Sending data...")

	defer wg.Done()

	bulkApiURL := baseApiUrl + bulkv2Url

	body := BulkDataRequest{
		IndexName: os.Getenv("INDEX_NAME"),
		Records:   records,
	}

	jsonData, err := utils.ConvertToJson(body)
	if err != nil {
		return fmt.Errorf("error converting to JSON: %v", err)
	} else {
		log.Printf("Convert to json successful...")
	}

	// Crear una solicitud POST
	req, err := http.NewRequest("POST", bulkApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating the request: %v", err)
	} else {
		log.Printf("Bulk data request successfully created...")
	}

	setBasicAuth(req)

	client := &http.Client{}

	//Make the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error when making the request: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Uploading data...")

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful bulk upload: %d", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusOK {
		log.Println("Data uploaded...")
	}

	return nil
}

/*Search for records indexed by match*/
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
		MaxResults: 550000,
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

/*Search for records indexed by match*/
func SearchDataJSON(term string, page int, itemsPerPage int, searchType string) (models.EmailSearchResponse, error) {
	log.Println("Looking for information...")

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
		MaxResults: 500000,
		/*From:       page,
		MaxResults: itemsPerPage,*/
		//Source:     []
	}

	jsonData, err := utils.ConvertToJson(body)
	if err != nil {
		return models.EmailSearchResponse{}, err
	}

	req, err := http.NewRequest("POST", searchApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return models.EmailSearchResponse{}, fmt.Errorf("error al crear la solicitud: %v", err)
	}

	setBasicAuth(req)

	client := &http.Client{}

	// Realizar la solicitud POST
	resp, err := client.Do(req)
	if err != nil {
		return models.EmailSearchResponse{}, fmt.Errorf("error al realizar la solicitud: %v", err)
	}

	defer resp.Body.Close()

	// Leer la respuesta de la API
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EmailSearchResponse{}, fmt.Errorf("error al leer la respuesta de la API: %v", err)
	}

	var searchDataResponse SearchDataResponse
	if err := json.Unmarshal(responseBody, &searchDataResponse); err != nil {
		return models.EmailSearchResponse{}, fmt.Errorf("error al analizar la respuesta JSON: %v", err)
	}

	emails := make([]models.Email, 0)

	for _, hit := range searchDataResponse.Hits.Hits {
		emails = append(emails, hit.Source)
	}

	var result models.EmailSearchResponse
	result.TotalItems = len(emails)
	result.ItemsPerPage = itemsPerPage
	result.Page = page
	result.Items = emails

	return result, nil
}

/*Check if the index to be created already exists*/
func CheckIndexExists() (bool, error) {
	indexApiURL := baseApiUrl + existIndexUrl + utils.IndexName

	req, err := http.NewRequest("GET", indexApiURL, nil)
	if err != nil {
		return false, fmt.Errorf("error al crear la solicitud: %v", err)
	}

	setBasicAuth(req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error al realizar la solicitud: %v", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error al leer la respuesta de la API: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		log.Println("Index with name " + utils.IndexName + " exists...")
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		log.Println("Index with name " + utils.IndexName + " does not exist...")
		return false, nil
	}

	return false, fmt.Errorf("código de estado inesperado: %d", resp.StatusCode)
}
