package models

type ProgramFlags struct {
	ExcludeDir  string
	ExcludeFile string
	ThreadCount int
	ContextSize int
	Root        string
	Pattern     string
}
