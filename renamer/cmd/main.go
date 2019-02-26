package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"dgrindlay.io/episode-namer/renamer"
)

func main() {
	loginRequest := renamer.LoginRequest{Apikey: "0629B785CE550C8D"}
	loginResponseSuccess := new(renamer.LoginResponse)
	loginResponseError := new(renamer.ErrorResponse)

	loginErr := renamer.MakeLoginRequest(loginRequest, loginResponseSuccess, loginResponseError)

	if loginErr != nil {
		log.Fatal(loginErr)
		return
	}

	token := loginResponseSuccess.Token

	episodeFiles := renamer.MapEpisodeFiles(os.Args[1])

	for series, files := range episodeFiles {
		fmt.Println("Search for series: " + series)

		searchResponseSuccess := new(renamer.SearchResponse)
		searchResponseError := new(renamer.ErrorResponse)

		searchErr := renamer.MakeSearchRequest(series, token, searchResponseSuccess,
			searchResponseError)

		if searchErr != nil {
			log.Fatal(searchErr)
			return
		}

		seriesList := renamer.CullSeriesData(searchResponseSuccess.Data)

		reader := bufio.NewReader(os.Stdin)
		seriesID := -1
		seriesName := ""

		for {
			fmt.Println("Select index from list or q to exit: ")
			for i, series := range seriesList {
				fmt.Printf("%v: %v\n", i, series.SeriesName)
			}
			fmt.Println("Number of choices: ", len(seriesList))

			input, _ := reader.ReadString('\n')
			input = strings.Trim(input, "\n\r")

			if index, err := strconv.Atoi(input); err == nil {
				if index < len(seriesList) {
					fmt.Println("You have selected: " + seriesList[index].SeriesName)
					seriesName = seriesList[index].SeriesName
					confirmationInput, _ := reader.ReadString('\n')
					confirmationInput = strings.Trim(confirmationInput, "\n\r")
					seriesID = seriesList[index].ID
					break
				}

				fmt.Println("Index out of bounds max: " + strconv.Itoa(len(seriesList)))
				continue
			} else {
				fmt.Println(err)
			}

			fmt.Println("Not a valid input")
		}

		var episodes []renamer.EpisodeDetails

		page := 1
		for true {
			episodeDetailsSuccess := new(renamer.EpisodeResponse)
			episodeDetailsError := new(renamer.ErrorResponse)

			fmt.Println("Search for page: " + strconv.Itoa(page))
			err := renamer.GetEpisodeDetails(seriesID, token, page, episodeDetailsSuccess, episodeDetailsError)

			if err != nil {
				log.Fatal(err)
				break
			}

			if episodeDetailsError.Error != "" {
				out, err := json.Marshal(episodeDetailsError)
				if err != nil {
					log.Fatal(err)
					break
				}

				fmt.Println(string(out))
				break
			}

			fmt.Println("Number of episodes found: " + strconv.Itoa(len(episodeDetailsSuccess.Data)))

			for _, episode := range episodeDetailsSuccess.Data {
				episodes = append(episodes, episode)
			}

			if page == episodeDetailsSuccess.Links.Last {
				break
			}

			page++
		}

		fmt.Println("Total number of episodes: " + strconv.Itoa(len(episodes)))

		episodeNameMap := renamer.GetEpisodeNameMap(episodes)
		for k, v := range episodeNameMap {
			fmt.Printf("S%vE%v: %v\n", k.Season, k.Episode, v)
		}

		reader = bufio.NewReader(os.Stdin)
		filesToRename := renamer.RenameFiles(seriesName, files, episodeNameMap)
		for _, rename := range filesToRename {
			fmt.Printf("Rename file from %v to %v, y?\n", rename.From, rename.To)

			input, _ := reader.ReadString('\n')
			input = strings.Trim(input, "\n\r")

			if input == "y" {
				fmt.Println("Original: " + os.Args[1] + "/" + rename.From)
				fmt.Println("New: " + os.Args[1] + "/" + seriesName + "/Season " + rename.Season + "/" + rename.To)

				err := os.MkdirAll(os.Args[1]+"/"+seriesName+"/Season "+rename.Season, 0700)
				if err != nil {
					log.Fatal(err)
				}

				pathErr := os.Rename(os.Args[1]+"/"+rename.From, os.Args[1]+"/"+seriesName+"/Season "+rename.Season+"/"+rename.To)
				if err != nil {
					log.Fatal(pathErr)
				}
			}
		}
	}
}
