package models

type Email struct {
	MessageId               string `json:"messageId"`
	Date                    string `json:"date"`
	From                    string `json:"from"`
	To                      string `json:"to"`
	Subject                 string `json:"subject"`
	MimeVersion             string `json:"mimeVersion"`
	ContentType             string `json:"contentType"`
	ContentTransferEncoding string `json:"contentTransferEncoding"`
	XFrom                   string `json:"xFrom"`
	XTo                     string `json:"xTo"`
	XCc                     string `json:"xCc"`
	XBcc                    string `json:"xBcc"`
	XFolder                 string `json:"xFolder"`
	XOrigin                 string `json:"xOrigin"`
	XFileName               string `json:"xFileName"`
	Content                 string `json:"content"`
}
