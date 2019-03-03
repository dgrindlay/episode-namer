package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	r "github.com/dgrindlay/renamer"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: directory")
	}

	episodeFiles := r.MapEpisodeFiles(os.Args[1])

	tvdb := new(r.Tvdb)
	loginRequest := r.LoginRequest{Apikey: "0629B785CE550C8D"}
	loginErr := tvdb.Login(loginRequest)

	if loginErr != nil {
		log.Fatal(loginErr)
		return
	}

	for series, files := range episodeFiles {
		fmt.Println("Search for series: " + series)

		search, err := tvdb.Search(series)

		if err != nil {
			log.Fatal(err)
			return
		}

		seriesList, checkErr := tvdb.OrderByPriority(search.Data)
		if checkErr != nil {
			fmt.Println(checkErr)
			break
		}

		reader := bufio.NewReader(os.Stdin)
		seriesID := -1
		seriesName := ""

		limit := 5
		for {
			fmt.Println("Type index to select, c to show more or q to exit: ")
			for i := 0; i < limit; i++ {
				fmt.Printf("%v: %v\n", i, seriesList[i].SeriesName)

				if i == len(seriesList)-1 {
					break
				}
			}

			input, _ := reader.ReadString('\n')
			input = strings.Trim(input, "\n\r")

			if input == "q" {
				fmt.Println("Exiting...")
				os.Exit(0)
			}

			if input == "c" {
				fmt.Println("Continuing...")
				limit += 5
				continue
			}

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
			}

			fmt.Println("Not a valid input")
		}

		var episodes []r.EpisodeDetails

		page := 1
		for true {
			fmt.Println("Search for page: " + strconv.Itoa(page))
			episodeResponse, err := tvdb.GetEpisodes(seriesID, page)

			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println("Number of episodes found: " + strconv.Itoa(len(episodeResponse.Data)))

			for _, episode := range episodeResponse.Data {
				episodes = append(episodes, episode)
			}

			if page == episodeResponse.Links.Last {
				break
			}

			page++
		}

		fmt.Println("Total number of episodes: " + strconv.Itoa(len(episodes)))

		episodeNameMap := r.GetEpisodeNameMap(episodes)
		for k, v := range episodeNameMap {
			fmt.Printf("S%vE%v: %v\n", k.Season, k.Episode, v)
		}

		reader = bufio.NewReader(os.Stdin)
		filesToRename := r.RenameFiles(seriesName, files, episodeNameMap)
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
