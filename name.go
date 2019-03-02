package renamer

import (
	"regexp"
	"strconv"
	"strings"
)

// GetEpisodeName returns a sanitized episode name for lookup.
func GetEpisodeName(fileName string) string {
	unwanted := []string{"web", "720p", "1080p", "REPACK", "WEBRip", "AAC", "x264"}
	for _, word := range unwanted {
		fileName = strings.Replace(fileName, word, "", -1)
	}

	episodeName := removeFileType(fileName)
	episodeName = removeBrackets(episodeName)

	episodeName = strings.Replace(episodeName, ".", " ", -1)
	episodeName = strings.Replace(episodeName, "_", " ", -1)
	episodeName = strings.Replace(episodeName, "  ", " ", -1)

	episodeName = extractName(episodeName)

	episodeName = removeEpisodeDetails(episodeName)
	episodeName = removeExtraWhiteSpace(episodeName)

	return strings.Trim(episodeName, " ")
}

func removeBrackets(name string) string {
	roundBracketRegex, _ := regexp.Compile("\\(.*?\\)")
	squareBracketRegex, _ := regexp.Compile("\\[.*?\\]")

	noRoundBrackets := roundBracketRegex.ReplaceAllString(name, "")
	noSquareBrackets := squareBracketRegex.ReplaceAllString(noRoundBrackets, "")
	return noSquareBrackets
}

func removeFileType(name string) string {
	lastIndex := strings.LastIndex(name, ".")

	if lastIndex != -1 {
		return name[:lastIndex]
	}

	return name
}

func GetFileType(name string) string {
	lastIndex := strings.LastIndex(name, ".")

	if lastIndex != -1 {
		return name[lastIndex:]
	}

	return name
}

// GetEpisodeNumber get the episode and season number from file name
func GetEpisodeNumber(name string) EpisodeNumber {
	explicitRegex, _ := regexp.Compile("[sS]\\d+[eE]\\d+")
	explicitSubString := explicitRegex.FindString(name)

	episodeRegex, _ := regexp.Compile("[eE]\\d+")
	seasonRegex, _ := regexp.Compile("[sS]\\d+")

	var episode = -1
	var season = -1

	if episodeRegex.MatchString(explicitSubString) && seasonRegex.MatchString(explicitSubString) {
		episode, _ = strconv.Atoi(episodeRegex.FindString(explicitSubString)[1:])
		season, _ = strconv.Atoi(seasonRegex.FindString(explicitSubString)[1:])
		return EpisodeNumber{season, episode}
	}

	implicitRegex, _ := regexp.Compile("\\d+x\\d+")
	implicitSubString := implicitRegex.FindString(name)

	parts := strings.Split(implicitSubString, "x")

	digitRegex, _ := regexp.Compile("\\d+")
	if digitRegex.MatchString(parts[0]) {
		season, _ = strconv.Atoi(parts[0])
	}

	if digitRegex.MatchString(parts[1]) {
		episode, _ = strconv.Atoi(parts[1])
	}

	return EpisodeNumber{season, episode}
}

func removeEpisodeDetails(name string) string {
	episodeRegex, _ := regexp.Compile("[sS]\\d+[eE]\\d+")
	return episodeRegex.ReplaceAllString(name, "")
}

func extractName(name string) string {
	var startingName []string

	digitRegex, _ := regexp.Compile("\\d")
	wordRegex, _ := regexp.Compile("\\w")

	words := strings.Split(name, " ")
	for _, word := range words {
		word = strings.Trim(word, "\\")
		if digitRegex.MatchString(word) && wordRegex.MatchString(word) {
			break
		} else {
			startingName = append(startingName, word)
		}
	}

	return strings.Join(startingName, " ")
}

func removeExtraWhiteSpace(name string) string {
	whiteSpace, _ := regexp.Compile("\\s+")
	return whiteSpace.ReplaceAllString(name, " ")
}
