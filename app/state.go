package app

import "github.com/fiifoo/mp3rand/fs"

type State struct {
	Phase int
	Done  bool
	Error error

	Source Source
	Target Target
}

type Source struct {
	Directory string
	TotalSize int64
	Files     []fs.File
}

type Target struct {
	Directory string
	MaxSize   int64
	Files     []fs.File
}
