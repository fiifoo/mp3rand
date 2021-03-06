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
		answer := ui.Query("Source directory")

		if len(answer) > 0 && fs.IsDirectory(answer) {
			state.Source.Directory = answer

			state.Phase++
		}
	},
	func(state *State) {
		state.Target = Target{}
		answer := ui.Query("Target directory")

		if len(answer) > 0 && fs.IsDirectory(answer) && answer != state.Source.Directory {
			state.Target.Directory = answer

			state.Phase++
		}
	},
	func(state *State) {
		available := fs.SpaceAvailable(state.Target.Directory)
		state.Target.MaxSize = available

		answer := ui.Query(fmt.Sprintf("Max size to copy in MB (enter = use all free space available = %d MB)", available/fs.MB))

		if answer == "" {
			state.Phase++
		} else {
			maxSize, err := strconv.Atoi(answer)

			if err == nil && maxSize > 0 {
				state.Target.MaxSize = int64(maxSize) * fs.MB

				state.Phase++
			}
		}
	},
	func(state *State) {
		fmt.Printf("Copying random files (max %d MB) from %s to %s\n\n", state.Target.MaxSize/fs.MB, state.Source.Directory, state.Target.Directory)

		if !fs.IsEmptyDirectory(state.Target.Directory) {
			fmt.Println("**WARNING** Target directory is not empty! Copy will be executed anyway if you confirm it **WARNING**")
			fmt.Println()
		}

		ui.Query("Press enter to confirm")
		state.Phase++
	},
	func(state *State) {
		progress := func(totalCount int, totalSize int64) {
			ui.Clear()
			fmt.Printf("Reading %s: total files %d, total size %d MB \n\n", state.Source.Directory, totalCount, totalSize/fs.MB)
		}

		files, totalSize, err := fs.ReadFiles(state.Source.Directory, fileExtension, progress)

		if err != nil {
			state.Error = err
		} else {
			if len(files) == 0 {
				state.Error = fmt.Errorf("No files found in source directory %s matching file extension %s", state.Source.Directory, fileExtension)
			} else {
				state.Source.Files = files
				state.Source.TotalSize = totalSize
				state.Phase++
			}
		}
	},
	func(state *State) {
		state.Target.Files = fs.SelectRandomFiles(state.Source.Files, state.Target.MaxSize)

		progress := func(totalCount int) {
			ui.Clear()
			fmt.Printf("Writing to %s: files %d / %d \n\n", state.Target.Directory, totalCount, len(state.Target.Files))
		}

		err := fs.CopyFilesTo(state.Target.Files, state.Target.Directory, progress)

		if err != nil {
			state.Error = err
		} else {
			state.Phase++
		}
	},
	func(state *State) {
		fmt.Printf("Ready! Copied %d files to %s\n\n", len(state.Target.Files), state.Target.Directory)

		ui.Query("Press enter to quit")
		state.Done = true
	},
}

func Run() {
	state := State{}

	for !state.Done {
		if state.Error != nil {
			fmt.Println("Error: " + state.Error.Error())
			fmt.Println()

			ui.Query("Press enter to quit")
			state.Done = true
		} else {
			phase := phases[state.Phase]

			phase(&state)
		}
	}
}
