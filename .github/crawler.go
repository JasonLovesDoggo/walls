package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
)

const (
	DescriptionFile     = "DESCRIPTION"
	MaxFilesPerCategory = 2
)

var (
	BaseDir            string
	IgnoredDirectories = []string{".github", ".idea", ".git"}
	IgnoredFiles       = []string{".gitignore", OutputFile, DescriptionFile, ".DS_Store", ".gitkeep"}
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.Dir(cwd)
}

func shouldIgnore(info os.FileInfo) bool {
	if slices.Contains(IgnoredDirectories, info.Name()) || slices.Contains(IgnoredFiles,
		info.Name()) {
		return true
	}
	return false
}

func getEstimatedCollectionsCount(root string) int {
	var count int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if shouldIgnore(info) {
			return filepath.SkipDir
		}

		count++
		return nil
	})
	if err != nil {
		panic(err)
	}
	return count
}

func ParseDirectories(root string) []Collection {

	collections := make([]Collection, 0, getEstimatedCollectionsCount(root))
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldIgnore(info) || info.Name() == BaseDir {
			return filepath.SkipDir
		}

		collection, err := ParseDirectory(path)
		if err != nil {
			return err
		}
		collections = append(collections, collection)
		return nil
	})
	if err != nil {
		return nil
	}
	return collections
}

func ParseDirectory(path string) (Collection, error) {

	exists, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !exists.IsDir() {
		return Collection{}, fmt.Errorf("%s is not a directory", path)
	}

	files := getRandomizedFiles(path, MaxFilesPerCategory)

	return Collection{
		Name:        filepath.Base(path),
		Files:       files,
		BrowseUrl:   "./" + filepath.Base(path),
		Description: GetDescription(path),
	}, nil
}

// getRandomizedFiles returns a list of files in the directory in random order
func getRandomizedFiles(path string, limit uint) []File {
	files, err := os.ReadDir(path)
	fmt.Println(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
	// Exclude files that are ignored
	for i, file := range files {
		stat, err := os.Stat(path + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		if shouldIgnore(stat) {
			files = slices.Delete(files, i, i+1)
			break
		}
	}
	fmt.Println(files)

	// Randomize the order of the files
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	// Limit the number of files
	if limit > 0 && uint(len(files)) > limit {
		files = files[:limit]
	}

	fmt.Println(files)
	var result []File
	for _, file := range files {
		result = append(result, File{
			Url:         "./" + file.Name(),
			Description: file.Name(),
		})
	}

	return result

}

// GetDescription returns the description of the directory, assumes the directory has a DescriptionFile in it
func GetDescription(path string) string {
	if _, err := os.Stat(filepath.Join(path, DescriptionFile)); err == nil {
		file, err := os.ReadFile(filepath.Join(path, DescriptionFile))
		if err != nil {
			return ""
		}
		return string(file)
	}

	return ""

}
