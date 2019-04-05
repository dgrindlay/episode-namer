package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dgrindlay/renamer"
)

func main() {
	animeList := new(renamer.AniList)
	searchResults, err := animeList.Search("my hero academia")

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type index to select, c to show more or q to exit: ")
	for i, result := range searchResults.Page.Media {
		fmt.Printf("%v: %v\n", i, result.Title.English)
	}

	input, _ := reader.ReadString('\n')
	input = strings.Trim(input, "\n\r")

	if input == "q" {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	if index, err := strconv.Atoi(input); err == nil {
		if index < len(searchResults.Page.Media) {
			fmt.Println("You have selected: " + searchResults.Page.Media[index].Title.English)
			seriesName = seriesList[index].SeriesName
			confirmationInput, _ := reader.ReadString('\n')
			confirmationInput = strings.Trim(confirmationInput, "\n\r")
			seriesID = seriesList[index].ID
			break
		}
	}

	fmt.Println("Not a valid input")
}
