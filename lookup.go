package renamer

import (
	"fmt"
	"sort"
	"strconv"
)

// GetEpisodeNameMap returns a map of season, episode number keys to the name of the episode
func GetEpisodeNameMap(episodeDetails []EpisodeDetails) map[EpisodeNumber]string {
	var episodeNumberMap map[EpisodeNumber]string
	episodeNumberMap = make(map[EpisodeNumber]string)

	for _, episode := range episodeDetails {
		number := EpisodeNumber{Season: episode.AiredSeason, Episode: episode.AiredEpisodeNumber}
		episodeNumberMap[number] = episode.EpisodeName
	}

	return episodeNumberMap
}

type byID []SeriesData

func (s byID) Len() int {
	return len(s)
}

func (s byID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byID) Less(i, j int) bool {
	return s[i].ID <= s[j].ID
}

// CullSeriesData returns an ordered list of the most likely series data. Removes data that do
// not have aliases, orders by firstAired and by if it is still continuing.
func CullSeriesData(seriesData []SeriesData) []SeriesData {
	sort.Sort(byID(seriesData))
	// sort.Sort(byAiredDate(orderedSeriesData))
	// sort.Sort(byStatus(orderedSeriesData))
	return seriesData
}

// RenameTuple contains all the information required to rename a file from it's original name
// to it's season, episode name
type RenameTuple struct {
	From   string
	To     string
	Season string
}

// RenameFiles renames files to their episode name
func RenameFiles(seriesName string, files []string, episodeMap map[EpisodeNumber]string) []RenameTuple {
	var filesToRename []RenameTuple

	for _, fileName := range files {
		episodeNumber := GetEpisodeNumber(fileName)

		if episodeNumber.Episode == -1 {
			fmt.Println("Could not find episode number for " + fileName)
			continue
		}

		episodeName, ok := episodeMap[episodeNumber]
		if ok {
			fileExtension := GetFileType(fileName)
			fullName := fmt.Sprintf("%v - s%02de%02d - %v%v", seriesName, episodeNumber.Season, episodeNumber.Episode, episodeName, fileExtension)

			fileRename := RenameTuple{From: fileName, To: fullName, Season: strconv.Itoa(episodeNumber.Season)}
			filesToRename = append(filesToRename, fileRename)
			continue
		}

		fmt.Printf("Could not find episode name for S%vE%v\n", episodeNumber.Season, episodeNumber.Episode)
	}

	return filesToRename
}
