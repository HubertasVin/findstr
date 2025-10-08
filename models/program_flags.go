package models

type ProgramFlags struct {
	ExcludeDir  string
	ExcludeFile string
	ThreadCount int
	ContextSize int
	Root        string
	SkipGit     bool
	SearchArch  bool
    Json        bool
	Pattern     string
}
