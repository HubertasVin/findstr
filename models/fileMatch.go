package models

type FileMatch struct {
    File            string
    ContextLineNums []int
    HighLineNums    []int
    FileContent     []string
}
