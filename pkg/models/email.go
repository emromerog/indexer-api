package models

type Email struct {
	//Id        int    `json:"id"`
	MessageId string `json:"messageId"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	/*MimeVersion             string `json:"mimeVersion"`
	ContentType             string `json:"contentType"`
	ContentTransferEncoding string `json:"contentTransferEncoding"`
	XFrom                   string `json:"xFrom"`
	XTo                     string `json:"xTo"`
	XCc                     string `json:"xCc"`
	XBcc                    string `json:"xBcc"`
	XFolder                 string `json:"xFolder"`
	XOrigin                 string `json:"xOrigin"`*/
	XFileName string `json:"xFileName"`
	Content   string `json:"content"`
}

type EmailSearchRequest struct {
	Term         string `json:"term"`
	Page         int    `json:"page"`
	ItemsPerPage int    `json:"maxResults"`
}

type EmailSearchResponse struct {
	TotalItems   int     `json:"totalItems"`
	Items        []Email `json:"items"`
	Page         int     `json:"page"`
	ItemsPerPage int     `json:"itemsPerPage"`
}
