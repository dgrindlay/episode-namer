package renamer

import (
	"errors"
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

type byRating struct {
	seriesData  []*SeriesData
	priorityMap map[*SeriesData]int
}

func (s byRating) Len() int {
	return len(s.seriesData)
}

func (s byRating) Swap(i, j int) {
	s.seriesData[i], s.seriesData[j] = s.seriesData[j], s.seriesData[i]
}

func (s byRating) Less(i, j int) bool {
	return s.priorityMap[s.seriesData[i]] >= s.priorityMap[s.seriesData[j]]
}

type SeriesPriority struct {
	series     *SeriesData
	priorities []float64
}

func normalizeSeriesPriorities(seriesPriority []SeriesPriority) ([]SeriesPriority, error) {
	normalizedSeriesPriorities := []SeriesPriority{}

	if len(seriesPriority) == 0 || len(seriesPriority[0].priorities) == 0 {
		return normalizedSeriesPriorities, errors.New("Length error")
	}

	maxPriorities := make([]float64, len(seriesPriority[0].priorities))

	for _, series := range seriesPriority {
		seriesPriorities := series.priorities
		if len(seriesPriorities) != len(maxPriorities) {
			return normalizedSeriesPriorities, errors.New("Length error")
		}

		for i, priority := range seriesPriorities {
			if priority > maxPriorities[i] {
				maxPriorities[i] = priority
			}
		}
	}

	for _, series := range seriesPriority {
		normalPriorities := make([]float64, len(maxPriorities))
		for i := 0; i < len(maxPriorities); i++ {
			if maxPriorities[i] != 0 {
				normalPriorities[i] = (series.priorities[i] / maxPriorities[i]) * 100
			} else {
				normalPriorities[i] = 0
			}
		}

		normalizedSeriesPriorities = append(normalizedSeriesPriorities, SeriesPriority{series.series, normalPriorities})
	}

	return normalizedSeriesPriorities, nil
}
