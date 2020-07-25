package app

import "github.com/fiifoo/mp3rand/fs"

type State struct {
	Phase int
	Done  bool

	Source Source
	Target Target
}

type Source struct {
	Directory string
	Files     []fs.File
	TotalSize int64
}

type Target struct {
	Directory string
	MaxSize   int64
	Files     []fs.File
}
