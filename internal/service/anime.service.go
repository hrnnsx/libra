package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/hrnnsx/libra/external/anilist"
	"github.com/hrnnsx/libra/external/tracemoe"
)

var ErrAnimeNotFound = errors.New("anime not found")

type AnimeService interface {
	BrowseAnime(
		ctx context.Context,
		page int,
		perPage int,
	) (*AnimeListResponse, error)

	SearchAnime(
		ctx context.Context,
		params AnimeSearchParams,
	) (*AnimeListResponse, error)

	GetAnime(
		ctx context.Context,
		externalID string,
	) (*AnimeResponse, error)

	IdentifyAnime(
		ctx context.Context,
		imageURL string,
	) (*AnimeIdentificationListResponse, error)
}

type AnimeSearchParams struct {
	Title      string
	ExternalID string
	Season     string
	SeasonYear int
	Genre      string
	Status     string
	Format     string
	Page       int
	PerPage    int
}

type AnimeResponse struct {
	ExternalID string     `json:"external_id"`
	Title      AnimeTitle `json:"title"`
	Synopsis   string     `json:"synopsis"`
	ImageURL   string     `json:"image_url"`
	Episodes   *int       `json:"episodes,omitempty"`
	Status     string     `json:"status,omitempty"`
	Format     string     `json:"format,omitempty"`
	Season     string     `json:"season,omitempty"`
	SeasonYear *int       `json:"season_year,omitempty"`
	Genres     []string   `json:"genres,omitempty"`
	Score      *int       `json:"score,omitempty"`
}

type AnimeTitle struct {
	Romaji  string `json:"romaji,omitempty"`
	English string `json:"english,omitempty"`
	Native  string `json:"native,omitempty"`
}

type AnimeListResponse struct {
	Data       []AnimeResponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type Pagination struct {
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	HasNextPage bool `json:"has_next_page"`
	Total       int  `json:"total"`
}

type AnimeIdentificationResponse struct {
	AnilistID  int     `json:"anilist_id"`
	Filename   string  `json:"filename"`
	Episode    *int    `json:"episode,omitempty"`
	From       float64 `json:"from"`
	To         float64 `json:"to"`
	Similarity float64 `json:"similarity"`
	Video      string  `json:"video"`
	Image      string  `json:"image"`
}

type AnimeIdentificationListResponse struct {
	Data []AnimeIdentificationResponse `json:"data"`
}

type animeService struct {
	anilistClient  *anilist.Client
	tracemoeClient *tracemoe.Client
}

func NewAnimeService(
	anilistClient *anilist.Client,
	tracemoeClient *tracemoe.Client,
) AnimeService {
	return &animeService{
		anilistClient:  anilistClient,
		tracemoeClient: tracemoeClient,
	}
}

func (s *animeService) BrowseAnime(
	ctx context.Context,
	page int,
	perPage int,
) (*AnimeListResponse, error) {
	result, err := s.anilistClient.BrowseAnime(
		ctx,
		page,
		perPage,
	)
	if err != nil {
		return nil, err
	}

	data := make([]AnimeResponse, 0, len(result.Anime))

	for _, anime := range result.Anime {
		data = append(data, mapAnime(anime))
	}

	return &AnimeListResponse{
		Data: data,
		Pagination: Pagination{
			Page:        result.PageInfo.CurrentPage,
			PerPage:     result.PageInfo.PerPage,
			HasNextPage: result.PageInfo.HasNextPage,
			Total:       result.PageInfo.Total,
		},
	}, nil
}

func (s *animeService) SearchAnime(
	ctx context.Context,
	params AnimeSearchParams,
) (*AnimeListResponse, error) {
	externalID := 0

	if params.ExternalID != "" {
		id, err := strconv.Atoi(params.ExternalID)
		if err != nil {
			return nil, err
		}

		externalID = id
	}

	result, err := s.anilistClient.SearchAnime(
		ctx,
		anilist.SearchParams{
			Title:      params.Title,
			ExternalID: externalID,
			Season:     params.Season,
			SeasonYear: params.SeasonYear,
			Genre:      params.Genre,
			Status:     params.Status,
			Format:     params.Format,
			Page:       params.Page,
			PerPage:    params.PerPage,
		},
	)
	if err != nil {
		return nil, err
	}

	data := make([]AnimeResponse, 0, len(result.Anime))

	for _, anime := range result.Anime {
		data = append(data, mapAnime(anime))
	}

	return &AnimeListResponse{
		Data: data,
		Pagination: Pagination{
			Page:        result.PageInfo.CurrentPage,
			PerPage:     result.PageInfo.PerPage,
			HasNextPage: result.PageInfo.HasNextPage,
			Total:       result.PageInfo.Total,
		},
	}, nil
}

func (s *animeService) GetAnime(
	ctx context.Context,
	externalID string,
) (*AnimeResponse, error) {
	id, err := strconv.Atoi(externalID)
	if err != nil {
		return nil, ErrAnimeNotFound
	}

	result, err := s.anilistClient.GetAnime(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapAnime(*result)

	return &response, nil
}

func (s *animeService) IdentifyAnime(
	ctx context.Context,
	imageURL string,
) (*AnimeIdentificationListResponse, error) {
	result, err := s.tracemoeClient.Search(
		ctx,
		imageURL,
	)
	if err != nil {
		return nil, err
	}

	data := make(
		[]AnimeIdentificationResponse,
		0,
		len(result.Result),
	)

	for _, item := range result.Result {
		data = append(data, AnimeIdentificationResponse{
			AnilistID:  item.AnilistID,
			Filename:   item.Filename,
			Episode:    item.Episode,
			From:       item.From,
			To:         item.To,
			Similarity: item.Similarity,
			Video:      item.Video,
			Image:      item.Image,
		})
	}

	return &AnimeIdentificationListResponse{
		Data: data,
	}, nil
}

func mapAnime(anime anilist.Anime) AnimeResponse {
	return AnimeResponse{
		ExternalID: strconv.Itoa(anime.ID),
		Title: AnimeTitle{
			Romaji:  anime.Title.Romaji,
			English: anime.Title.English,
			Native:  anime.Title.Native,
		},
		Synopsis:   anime.Description,
		ImageURL:   anime.CoverImage.Large,
		Episodes:   anime.Episodes,
		Status:     anime.Status,
		Format:     anime.Format,
		Season:     anime.Season,
		SeasonYear: anime.SeasonYear,
		Genres:     anime.Genres,
		Score:      anime.Score,
	}
}
