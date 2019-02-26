package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MapEpisodeFiles returns a map of series names to list
func MapEpisodeFiles(dirName string) map[string][]string {
	acceptedFormats := []string{".mkv", ".avi", ".mp4"}
	var videoFiles []string

	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		extensionIndex := strings.LastIndex(path, ".")

		if extensionIndex != -1 {
			fileExtension := path[extensionIndex:]
			for _, format := range acceptedFormats {
				if fileExtension == format {
					videoFiles = append(videoFiles, path)
				}
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	contentNameMap := make(map[string][]string)

	for _, videoFile := range videoFiles {
		fileNameIndex := strings.LastIndex(videoFile, "\\")
		fileName := videoFile[fileNameIndex:]
		fileName = strings.ToLower(fileName)
		fileName = strings.Trim(fileName, "\\")
		episodeName := GetEpisodeName(fileName)

		if contentNameMap[episodeName] != nil {
			contentNameMap[episodeName] = append(contentNameMap[episodeName], fileName)
		} else {
			contentNameMap[episodeName] = []string{fileName}
		}
	}

	for series, files := range contentNameMap {
		fmt.Println(series, files)
	}

	return contentNameMap
}
