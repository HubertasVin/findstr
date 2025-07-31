package models

type FileMatch struct {
	File            string
	ContextLineNums []int
	MatchLineNums   []int
	FileContent     []string
}
