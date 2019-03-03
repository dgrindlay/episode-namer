package renamer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Tvdb represents a class used to interact with the tvdb api
type Tvdb struct {
	token string
}

// Login make login request to tvdb api, returns JWT token or error
func (tvdb *Tvdb) Login(request LoginRequest) error {
	resp := new(LoginResponse)
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return err
	}

	loginResp, err := http.Post("https://api.thetvdb.com/login", "application/json", bytes.NewBuffer(jsonRequest))

	if err != nil {
		return err
	}

	defer loginResp.Body.Close()
	if loginResp.StatusCode == 200 {
		err := json.NewDecoder(loginResp.Body).Decode(resp)
		if err != nil {
			return err
		}

		tvdb.token = resp.Token
	} else if loginResp.StatusCode >= 400 {
		errMessage := new(ErrorResponse)
		err := json.NewDecoder(loginResp.Body).Decode(errMessage)
		if err != nil {
			return err
		}

		return fmt.Errorf("%v: %v", strconv.Itoa(loginResp.StatusCode), errMessage.Error)
	}

	return nil
}

// Search make search request to tvdb
func (tvdb *Tvdb) Search(searchTerm string) (*SearchResponse, error) {
	searchResp := new(SearchResponse)
	client := &http.Client{}

	searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
	req, err := http.NewRequest("GET", "https://api.thetvdb.com/search/series?name="+searchTerm, nil)

	if err != nil {
		fmt.Println(err)
		return searchResp, err
	}

	req.Header.Add("Authorization", "Bearer "+tvdb.token)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return searchResp, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		err := json.NewDecoder(resp.Body).Decode(searchResp)
		return searchResp, err
	} else if resp.StatusCode >= 400 {
		errMessage := new(ErrorResponse)
		err := json.NewDecoder(resp.Body).Decode(errMessage)

		if err != nil {
			return searchResp, err
		}

		return searchResp, fmt.Errorf("%v: %v", strconv.Itoa(resp.StatusCode), errMessage.Error)
	}

	return searchResp, errors.New("Unknown response")
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

func (tvdb *Tvdb) OrderByPriority(seriesData []SeriesData) ([]SeriesData, error) {
	orderedSeriesPointers := []*SeriesData{}
	orderedSeriesList := []SeriesData{}
	seriesPriorityMap := make(map[*SeriesData]int)
	client := &http.Client{}

	for i := range seriesData {
		series := seriesData[i]
		orderedSeriesPointers = append(orderedSeriesPointers, &series)

		req, err := http.NewRequest("GET", "https://api.thetvdb.com/series/"+strconv.Itoa(series.ID)+"/filter?keys=siteRating,siteRatingCount", nil)
		if err != nil {
			fmt.Println(err)
			return seriesData, err
		}

		req.Header.Add("Authorization", "Bearer "+tvdb.token)
		req.Header.Add("Accept", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			return seriesData, err
		}

		var data map[string]interface{}

		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			err := json.NewDecoder(resp.Body).Decode(&data)

			if err != nil {
				fmt.Println(err)
				return seriesData, err
			}

			respData := data["data"].(map[string]interface{})
			siteRating := respData["siteRating"].(float64)
			siteRatingCount := respData["siteRatingCount"].(float64)
			seriesPriorityMap[&series] = int((siteRating + siteRatingCount))
		} else {
			errMessage := new(ErrorResponse)
			err := json.NewDecoder(resp.Body).Decode(errMessage)

			if err != nil {
				return seriesData, err
			}

			return seriesData, fmt.Errorf("%v: %v", strconv.Itoa(resp.StatusCode), errMessage.Error)
		}
	}

	sort.Sort(byRating{seriesData: orderedSeriesPointers, priorityMap: seriesPriorityMap})
	for _, seriesPointer := range orderedSeriesPointers {
		orderedSeriesList = append(orderedSeriesList, *seriesPointer)
	}

	return orderedSeriesList, nil
}

// GetEpisodes makes a request to tvdb for epsiode details
func (tvdb *Tvdb) GetEpisodes(seriesID int, page int) (*EpisodeResponse, error) {
	episodeResp := new(EpisodeResponse)
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.thetvdb.com/series/"+strconv.Itoa(seriesID)+
		"/episodes?page="+strconv.Itoa(page), nil)

	if err != nil {
		fmt.Println(err)
		return episodeResp, err
	}

	req.Header.Add("Authorization", "Bearer "+tvdb.token)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return episodeResp, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		err := json.NewDecoder(resp.Body).Decode(episodeResp)
		return episodeResp, err
	} else if resp.StatusCode >= 400 {
		errMessage := new(ErrorResponse)
		err := json.NewDecoder(resp.Body).Decode(errMessage)

		if err != nil {
			return episodeResp, fmt.Errorf("%v: %v", strconv.Itoa(resp.StatusCode), errMessage.Error)
		}

		return episodeResp, err
	}

	return episodeResp, errors.New("Unknown response")
}
