package zincsearch

import (
	"github.com/emromerog/indexer-api/pkg/models"
)

type BulkDataRequest struct {
	IndexName string         `json:"index"`
	Records   []models.Email `json:"records"`
}

type SearchDataRequest struct {
	SearchType string      `json:"search_type"`
	Query      SearchQuery `json:"query"`
	From       int         `json:"from"`
	MaxResults int         `json:"max_results"`
	Source     []string    `json:"_source"`
}

type SearchQuery struct {
	Term      string `json:"term"`
	Field     string `json:"field"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
