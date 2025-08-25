package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
)

const (
	pageSize = 100
)

type Repo struct {
	Name     string
	Owner    string
	CloneURL string
}

type jsonRepo struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	CloneURL string `json:"clone_url"`
}

type Reader struct {
}

func (r Reader) AllRepos(ctx context.Context, token, affiliation string) iter.Seq2[Repo, error] {
	return func(yield func(Repo, error) bool) {
		for page := 1; ; page++ {
			repos, err := readPage(ctx, affiliation, token, page)

			if err != nil {
				yield(Repo{}, err)
				return
			}

			for _, repo := range repos {
				if !yield(Repo{Name: repo.Name, Owner: repo.Owner.Login, CloneURL: repo.CloneURL}, nil) {
					return
				}
			}

			if len(repos) == 0 {
				return
			}
		}
	}
}

func readPage(ctx context.Context, affiliation string, token string, page int) ([]jsonRepo, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://api.github.com/user/repos?affiliation=%v&per_page=%v&page=%v", affiliation, pageSize, page)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return unmarshal(res)
}

func unmarshal(res *http.Response) ([]jsonRepo, error) {
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v (%v)", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	readRepos := []jsonRepo{}
	if err = json.Unmarshal(body, &readRepos); err != nil {
		return nil, err
	}

	return readRepos, nil
}
