package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const endpoint = "https://graphql.anilist.co"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphQLError  `json:"errors,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

func (c *Client) doRequest(
	ctx context.Context,
	query string,
	variables map[string]interface{},
	result interface{},
) error {
	payload := graphQLRequest{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("anilist returned status %d", resp.StatusCode)
	}

	var response graphQLResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if len(response.Errors) > 0 {
		return fmt.Errorf(
			"anilist graphql error: %s",
			response.Errors[0].Message,
		)
	}

	if err := json.Unmarshal(response.Data, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) BrowseAnime(
	ctx context.Context,
	page int,
	perPage int,
) (*AnimePage, error) {
	var response struct {
		Page AnimePage `json:"Page"`
	}

	variables := map[string]interface{}{
		"page":    page,
		"perPage": perPage,
	}

	if err := c.doRequest(
		ctx,
		browseAnimeQuery,
		variables,
		&response,
	); err != nil {
		return nil, err
	}

	return &response.Page, nil
}

func (c *Client) SearchAnime(
	ctx context.Context,
	params SearchParams,
) (*AnimePage, error) {
	var response struct {
		Page AnimePage `json:"Page"`
	}

	variables := map[string]interface{}{
		"page":    params.Page,
		"perPage": params.PerPage,
	}

	if params.Title != "" {
		variables["search"] = params.Title
	}

	if params.ExternalID > 0 {
		variables["id"] = params.ExternalID
	}

	if params.Season != "" {
		variables["season"] = params.Season
	}

	if params.SeasonYear > 0 {
		variables["seasonYear"] = params.SeasonYear
	}

	if params.Genre != "" {
		variables["genre"] = params.Genre
	}

	if params.Status != "" {
		variables["status"] = params.Status
	}

	if params.Format != "" {
		variables["format"] = params.Format
	}

	if err := c.doRequest(
		ctx,
		searchAnimeQuery,
		variables,
		&response,
	); err != nil {
		return nil, err
	}

	return &response.Page, nil
}

func (c *Client) GetAnime(
	ctx context.Context,
	id int,
) (*Anime, error) {
	var response struct {
		Media *Anime `json:"Media"`
	}

	variables := map[string]interface{}{
		"id": id,
	}

	if err := c.doRequest(
		ctx,
		animeDetailQuery,
		variables,
		&response,
	); err != nil {
		return nil, err
	}

	if response.Media == nil {
		return nil, errors.New("anime not found")
	}

	return response.Media, nil
}
