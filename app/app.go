package app

import (
	"fmt"
	"strconv"

	"github.com/fiifoo/mp3rand/fs"
	"github.com/fiifoo/mp3rand/ui"
)

type Phase func(state *State)

const fileExtension = "mp3"

var phases = []Phase{
	func(state *State) {
		state.Source = Source{}
		directory := ui.Query("Source directory")

		if len(directory) > 0 && fs.IsDirectory(directory) {
			state.Source.Directory = directory

			progress := func(totalCount int, totalSize int64) {
				ui.Clear()
				fmt.Printf("Reading %s: total files %d, total size %d MB \n\n", state.Source.Directory, totalCount, totalSize/1000000)
			}

			files, totalSize, err := fs.ReadFiles(state.Source.Directory, fileExtension, progress)

			if err != nil {
				fmt.Println(err)
			} else {
				state.Source.Files = files
				state.Source.TotalSize = totalSize
				state.Phase++
			}

		}
	},
	func(state *State) {
		state.Target = Target{}
		directory := ui.Query("Target directory")

		if len(directory) > 0 && fs.IsDirectory(directory) {
			state.Target.Directory = directory

			state.Phase++
		}
	},
	func(state *State) {
		answer := ui.Query("Max size to copy in MB")
		maxSize, err := strconv.Atoi(answer)

		if err == nil && maxSize > 0 {
			state.Target.MaxSize = int64(maxSize) * 1000000
			state.Target.Files = fs.SelectRandomFiles(state.Source.Files, state.Target.MaxSize)

			state.Phase++
		}
	},
	func(state *State) {
		fmt.Printf("Copying selected %d random files (max %d MB) from %s to %s\n\n", len(state.Target.Files), state.Target.MaxSize/1000000, state.Source.Directory, state.Target.Directory)

		if !fs.IsEmptyDirectory(state.Target.Directory) {
			fmt.Println("**WARNING** Target directory is not empty! Copy will be executed anyway if you confirm it **WARNING**\n")
		}

		ui.Query("Press enter to confirm")

		progress := func(totalCount int) {
			ui.Clear()
			fmt.Printf("Writing to %s: total files %d / %d \n\n", state.Target.Directory, totalCount, len(state.Target.Files))
		}

		err := fs.CopyFilesTo(state.Target.Files, state.Target.Directory, progress)

		if err != nil {
			fmt.Println(err)
		} else {
			state.Phase++
		}
	},
	func(state *State) {
		fmt.Println("Ready!\n")

		ui.Query("Press enter to quit")
		state.Done = true
	},
}

func Run() {
	state := State{}

	for !state.Done {
		phase := phases[state.Phase]

		phase(&state)
	}
}
