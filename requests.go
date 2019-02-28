package renamer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// default key: 0629B785CE550C8D

// MakeLoginRequest make login request to tvdb api, returns JWT token or error
func MakeLoginRequest(request LoginRequest, successMap interface{}, errorMap interface{}) error {
	jsonString, jsonErr := json.Marshal(request)

	if jsonErr != nil {
		return jsonErr
	}

	resp, err := http.Post("https://api.thetvdb.com/login", "application/json", bytes.NewBuffer(jsonString))

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return json.NewDecoder(resp.Body).Decode(successMap)
	}

	return json.NewDecoder(resp.Body).Decode(errorMap)
}

// MakeSearchRequest make search request to tvdb
func MakeSearchRequest(searchTerm string, token string, successMap interface{}, errorMap interface{}) error {
	client := &http.Client{}

	searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
	req, err := http.NewRequest("GET", "https://api.thetvdb.com/search/series?name="+searchTerm, nil)

	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return json.NewDecoder(resp.Body).Decode(successMap)
	}

	return json.NewDecoder(resp.Body).Decode(errorMap)
}

// GetEpisodeDetails makes a request to tvdb for epsiode details
func GetEpisodeDetails(seriesID int, token string, page int, successMap interface{}, errorMap interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.thetvdb.com/series/"+strconv.Itoa(seriesID)+
		"/episodes?page="+strconv.Itoa(page), nil)

	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return json.NewDecoder(resp.Body).Decode(successMap)
	}

	return json.NewDecoder(resp.Body).Decode(errorMap)
}
