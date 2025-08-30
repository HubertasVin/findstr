package models

type ProgramFlags struct {
	ExcludeDir  string
	ExcludeFile string
	ThreadCount int
	ContextSize int
	Root        string
    Style       string
    Json        bool
	Pattern     string
}
