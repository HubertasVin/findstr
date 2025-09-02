package models

type JsonFileMatch struct {
	FileName       string        `json:"fileName"`
	MatchedContent []LineContent `json:"matchedContent"`
}

type LineContent struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
}
