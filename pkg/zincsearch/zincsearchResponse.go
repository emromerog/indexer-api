package zincsearch

import (
	"github.com/emromerog/indexer-api/pkg/models"
)

type SearchDataResponse struct {
	Took     int        `json:"took"`
	TimedOut bool       `json:"timed_out"`
	Shards   ShardsInfo `json:"_shards"`
	Hits     HitsInfo   `json:"hits"`
}

type ShardsInfo struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsInfo struct {
	Total    HitTotal    `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []EmailInfo `json:"hits"`
}

type HitTotal struct {
	Value int `json:"value"`
}

type EmailInfo struct {
	Index  string       `json:"_index"`
	Type   string       `json:"_type"`
	ID     string       `json:"_id"`
	Score  float64      `json:"_score"`
	Source models.Email `json:"_source"`
}

/*type EmailData struct {
	Timestamp time.Time `json:"@timestamp"`
	Content   string    `json:"content"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	To        string    `json:"to"`
}*/
