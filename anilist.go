package renamer

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

type AniList struct {
}

type searchResponse struct {
	Page struct {
		PageInfo struct {
			Total       int
			CurrentPage int
			LastPage    int
			HasNextPage bool
			PerPage     int
		}

		Media []struct {
			ID    int
			Title struct {
				English string
				Native  string
			}
		}
	}
}

func (anilist *AniList) Search(searchTerm string) (searchResponse, error) {

	client := graphql.NewClient("https://graphql.anilist.co")

	req := graphql.NewRequest(`
		query ($page: Int, $perPage: Int, $search: String) {
			Page (page: $page, perPage: $perPage) {
				pageInfo {
					total
					currentPage
					lastPage
					hasNextPage
					perPage
				}
				media (search: $search, type: ANIME, sort: POPULARITY_DESC) {
					id
					title {
						english
					}
				}
			}
		}
	`)

	req.Var("search", searchTerm)
	req.Var("page", 1)
	req.Var("perPage", 10)

	var res searchResponse
	if err := client.Run(context.Background(), req, &res); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return res, nil
}

func (anilist *AniList) GetEpisodes(mediaId int) error {
	client := graphql.NewClient("https://graphql.anilist.co")

	req := graphql.NewRequest(`
		query {
			Media (id: 16498) {
				title {
					english
				}

				streamingEpisodes {
					title
				}
			}
		}
	`)

	req.Var("mediaId", 16498)

	type response struct {
		Media struct {
			Title struct {
				English string
			}

			StreamingEpisodes []struct {
				Title string
			}
		}
	}

	var res response
	if err := client.Run(context.Background(), req, &res); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return nil
}
