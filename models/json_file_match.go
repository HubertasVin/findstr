package models

type JsonFileMatch struct {
	FileName       string
	MatchedContent []LineContent
}

type LineContent struct {
	LineNumber int
	Content    string
}
