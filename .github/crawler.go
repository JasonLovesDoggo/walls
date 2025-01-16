package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"text/template"
)

const (
	DescriptionFile     = "DESCRIPTION"
	MaxFilesPerCategory = 2
)

var (
	BaseDir                string
	IgnoredDirectories     = []string{".github", ".idea", ".git"}
	IgnoredFiles           = []string{".gitignore", OutputFile, DescriptionFile, ".DS_Store", ".gitkeep"}
	ImageExtensions        = []string{".jpg", ".jpeg", ".png", ".webp", ".mp4"}
	CategoryReadmeTemplate string
	RootReadmeTemplate     string
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.Dir(cwd)

	CategoryReadmeTemplateRaw, _ := os.ReadFile(CategoryTemplateFile)
	CategoryReadmeTemplate = string(CategoryReadmeTemplateRaw)

	RootReadmeTemplateRaw, _ := os.ReadFile(RootTemplateFile)
	RootReadmeTemplate = string(RootReadmeTemplateRaw)

}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return slices.Contains(ImageExtensions, ext)
}

func shouldIgnore(info os.FileInfo) bool {
	if slices.Contains(IgnoredDirectories, info.Name()) ||
		slices.Contains(IgnoredFiles, info.Name()) {
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
		if info.IsDir() && !slices.Contains(IgnoredDirectories, info.Name()) {
			count++
		}
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

		// Skip ignored directories and base directory
		if shouldIgnore(info) || path == root || info.Name() == filepath.Base(root) {
			if path != root && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process directories
		if !info.IsDir() {
			return nil
		}

		collection, err := ParseDirectory(path)
		if err != nil {
			fmt.Printf("Error parsing directory %s: %v\n", path, err)
			return nil
		}

		// Only add collections with images
		if len(collection.Files) > 0 {
			collections = append(collections, collection)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directories: %v\n", err)
		return nil
	}
	// sort the collection files by alphabetical order
	sort.SliceStable(collections, func(i, j int) bool {
		return collections[i].Name < collections[j].Name
	})

	return collections
}

func ParseDirectory(path string) (Collection, error) {
	exists, err := os.Stat(path)
	if err != nil {
		return Collection{}, err
	}
	if !exists.IsDir() {
		return Collection{}, fmt.Errorf("%s is not a directory", path)
	}

	// Collect image files
	imageFiles, err := collectImageFiles(path)
	if err != nil {
		return Collection{}, err
	}

	collectionImageFiles := slices.Clone(imageFiles)
	// remove the ./dirname prefix from all the image urls
	for i := range collectionImageFiles {
		collectionImageFiles[i].Url = strings.Replace(collectionImageFiles[i].Url, "./"+filepath.Base(path)+"/", "./", 1)
	}

	// Create category README
	err = createCategoryReadme(path, Collection{
		Name:        filepath.Base(path),
		Files:       collectionImageFiles,
		BrowseUrl:   "./" + filepath.Base(path),
		Description: GetDescription(path),
	})
	if err != nil {
		fmt.Printf("Error creating category README for %s: %v\n", path, err)
	}

	return Collection{
		Name:        filepath.Base(path),
		Files:       getRandomizedFiles(imageFiles, MaxFilesPerCategory),
		BrowseUrl:   "./" + filepath.Base(path),
		Description: GetDescription(path),
	}, nil
}

func collectImageFiles(path string) ([]File, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var imageFiles []File
	for _, file := range files {
		// Skip directories and ignored files
		if file.IsDir() || shouldIgnore(fileInfoFromDirEntry(file)) {
			continue
		}

		// Check if it's an image file
		if isImage(file.Name()) {
			imageFiles = append(imageFiles, File{
				Url:         "./" + filepath.Base(path) + "/" + file.Name(),
				Description: file.Name(),
			})
		}
	}

	return imageFiles, nil
}

func getRandomizedFiles(files []File, limit uint) []File {
	// Randomize the order of the files
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})

	// Limit the number of files
	if limit > 0 && uint(len(files)) > limit {
		files = files[:limit]
	}

	return files
}

func createCategoryReadme(path string, collection Collection) error {
	// Create category README file
	readmePath := filepath.Join(path, "README.md")

	// Parse the category README template
	tmpl, err := template.New("CategoryReadme").Parse(CategoryReadmeTemplate)
	if err != nil {
		return err
	}

	// Create the README file
	f, err := os.Create(readmePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute the template
	err = tmpl.Execute(f, collection)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to convert os.DirEntry to os.FileInfo
func fileInfoFromDirEntry(entry os.DirEntry) os.FileInfo {
	info, err := entry.Info()
	if err != nil {
		panic(err)
	}
	return info
}

// GetDescription returns the description of the directory, assumes the directory has a DescriptionFile in it
func GetDescription(path string) string {
	descPath := filepath.Join(path, DescriptionFile)
	if _, err := os.Stat(descPath); err == nil {
		file, err := os.ReadFile(descPath)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(file))
	}

	return ""
}
