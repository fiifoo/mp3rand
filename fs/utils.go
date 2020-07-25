package fs

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type File struct {
	Path string
	Name string
	Size int64
}

func IsDirectory(directory string) bool {
	file, err := os.Open(directory)

	if err != nil {
		return false
	}

	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		return false
	}

	return info.IsDir()
}

func IsEmptyDirectory(directory string) bool {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return false
	}

	return len(files) == 0
}

func ReadFiles(
	directory string,
	extension string,
	progress func(totalCount int, totalSize int64)) (files []File, totalSize int64, err error) {

	files = make([]File, 0)
	totalCount := 0

	err = filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && hasExtension(path, extension) {
				file := File{path, info.Name(), info.Size()}

				totalCount++
				totalSize = totalSize + file.Size
				files = append(files, file)

				progress(totalCount, totalSize)
			}

			return nil
		})

	return
}

func SelectRandomFiles(files []File, maxSize int64) (result []File) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	var totalSize int64 = 0

	for _, file := range files {
		totalSize = totalSize + file.Size

		if totalSize <= maxSize {
			result = append(result, file)
		} else {
			break
		}
	}

	return
}

func CopyFilesTo(files []File, directory string, progress func(totalCount int)) (err error) {
	for i, file := range files {
		err = copyFileTo(file, directory)

		if err != nil {
			break
		}

		progress(i + 1)
	}

	return
}

func copyFileTo(sourceFile File, directory string) (err error) {
	destinationPath := directory + "\\" + sourceFile.Name

	input, err := ioutil.ReadFile(sourceFile.Path)

	if err == nil {
		err = ioutil.WriteFile(destinationPath, input, 0644)
	}

	return
}

func hasExtension(file string, extension string) bool {
	match, _ := regexp.MatchString(`\.`+extension+`$`, file)

	return match
}
