package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hrnnsx/libra/external/anilist"
	"github.com/hrnnsx/libra/external/tracemoe"
)

type fakeAniListClient struct {
	browseResult *anilist.AnimePage
	browseErr    error

	searchResult *anilist.AnimePage
	searchErr    error

	getResult *anilist.Anime
	getErr    error
}

func (f *fakeAniListClient) BrowseAnime(
	ctx context.Context,
	page int,
	perPage int,
) (*anilist.AnimePage, error) {
	return f.browseResult, f.browseErr
}

func (f *fakeAniListClient) SearchAnime(
	ctx context.Context,
	params anilist.SearchParams,
) (*anilist.AnimePage, error) {
	return f.searchResult, f.searchErr
}

func (f *fakeAniListClient) GetAnime(
	ctx context.Context,
	id int,
) (*anilist.Anime, error) {
	return f.getResult, f.getErr
}

type fakeTraceMoeClient struct {
	result *tracemoe.SearchResponse
	err    error
}

func (f *fakeTraceMoeClient) Search(
	ctx context.Context,
	imageURL string,
) (*tracemoe.SearchResponse, error) {
	return f.result, f.err
}

func TestAnimeService_BrowseAnime(t *testing.T) {
	client := &fakeAniListClient{
		browseResult: &anilist.AnimePage{
			PageInfo: anilist.PageInfo{
				CurrentPage: 1,
				HasNextPage: true,
				PerPage:     2,
				Total:       10,
			},
			Anime: []anilist.Anime{
				{
					ID: 20464,
					Title: anilist.AnimeTitle{
						Romaji:  "Mahouka Koukou no Rettousei",
						English: "The Irregular at Magic High School",
					},
					Description: "Test synopsis",
					CoverImage: anilist.CoverImage{
						Large: "https://example.com/image.jpg",
					},
				},
				{
					ID: 101,
					Title: anilist.AnimeTitle{
						Romaji: "Test Anime",
					},
				},
			},
		},
	}

	service := NewAnimeService(
		client,
		&fakeTraceMoeClient{},
	)

	result, err := service.BrowseAnime(
		context.Background(),
		1,
		2,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Data) != 2 {
		t.Fatalf(
			"expected 2 anime, got %d",
			len(result.Data),
		)
	}

	if result.Data[0].ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			result.Data[0].ExternalID,
		)
	}

	if result.Data[0].Title.Romaji != "Mahouka Koukou no Rettousei" {
		t.Errorf(
			"unexpected title: %s",
			result.Data[0].Title.Romaji,
		)
	}

	if result.Pagination.Page != 1 {
		t.Errorf(
			"expected page 1, got %d",
			result.Pagination.Page,
		)
	}

	if result.Pagination.PerPage != 2 {
		t.Errorf(
			"expected per_page 2, got %d",
			result.Pagination.PerPage,
		)
	}

	if !result.Pagination.HasNextPage {
		t.Error("expected has_next_page to be true")
	}

	if result.Pagination.Total != 10 {
		t.Errorf(
			"expected total 10, got %d",
			result.Pagination.Total,
		)
	}
}

func TestAnimeService_BrowseAnime_Error(t *testing.T) {
	expectedErr := errors.New("anilist unavailable")

	client := &fakeAniListClient{
		browseErr: expectedErr,
	}

	service := NewAnimeService(
		client,
		&fakeTraceMoeClient{},
	)

	_, err := service.BrowseAnime(
		context.Background(),
		1,
		20,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestAnimeService_SearchAnime(t *testing.T) {
	client := &fakeAniListClient{
		searchResult: &anilist.AnimePage{
			PageInfo: anilist.PageInfo{
				CurrentPage: 1,
				HasNextPage: false,
				PerPage:     1,
				Total:       1,
			},
			Anime: []anilist.Anime{
				{
					ID: 20464,
					Title: anilist.AnimeTitle{
						Romaji: "Mahouka Koukou no Rettousei",
					},
				},
			},
		},
	}

	service := NewAnimeService(
		client,
		&fakeTraceMoeClient{},
	)

	result, err := service.SearchAnime(
		context.Background(),
		AnimeSearchParams{
			Title:      "mahouka",
			ExternalID: "20464",
			Season:     "SPRING",
			SeasonYear: 2014,
			Genre:      "Action",
			Status:     "FINISHED",
			Format:     "TV",
			Page:       1,
			PerPage:    1,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf(
			"expected 1 anime, got %d",
			len(result.Data),
		)
	}

	if result.Data[0].ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			result.Data[0].ExternalID,
		)
	}
}

func TestAnimeService_SearchAnime_InvalidExternalID(t *testing.T) {
	service := NewAnimeService(
		&fakeAniListClient{},
		&fakeTraceMoeClient{},
	)

	_, err := service.SearchAnime(
		context.Background(),
		AnimeSearchParams{
			ExternalID: "abc",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAnimeService_GetAnime(t *testing.T) {
	client := &fakeAniListClient{
		getResult: &anilist.Anime{
			ID: 20464,
			Title: anilist.AnimeTitle{
				Romaji:  "Mahouka Koukou no Rettousei",
				English: "The Irregular at Magic High School",
			},
			Description: "Test synopsis",
			CoverImage: anilist.CoverImage{
				Large: "https://example.com/image.jpg",
			},
			Episodes: intPtr(26),
			Status:   "FINISHED",
			Format:   "TV",
		},
	}

	service := NewAnimeService(
		client,
		&fakeTraceMoeClient{},
	)

	result, err := service.GetAnime(
		context.Background(),
		"20464",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			result.ExternalID,
		)
	}

	if result.Title.Romaji != "Mahouka Koukou no Rettousei" {
		t.Errorf(
			"unexpected title: %s",
			result.Title.Romaji,
		)
	}

	if result.Episodes == nil || *result.Episodes != 26 {
		t.Errorf("expected 26 episodes")
	}

	if result.Status != "FINISHED" {
		t.Errorf(
			"expected FINISHED, got %s",
			result.Status,
		)
	}
}

func TestAnimeService_GetAnime_InvalidExternalID(t *testing.T) {
	service := NewAnimeService(
		&fakeAniListClient{},
		&fakeTraceMoeClient{},
	)

	_, err := service.GetAnime(
		context.Background(),
		"abc",
	)

	if !errors.Is(err, ErrAnimeNotFound) {
		t.Fatalf(
			"expected ErrAnimeNotFound, got %v",
			err,
		)
	}
}

func TestAnimeService_GetAnime_NotFound(t *testing.T) {
	client := &fakeAniListClient{
		getErr: ErrAnimeNotFound,
	}

	service := NewAnimeService(
		client,
		&fakeTraceMoeClient{},
	)

	_, err := service.GetAnime(
		context.Background(),
		"999999",
	)

	if !errors.Is(err, ErrAnimeNotFound) {
		t.Fatalf(
			"expected ErrAnimeNotFound, got %v",
			err,
		)
	}
}

func TestAnimeService_IdentifyAnime(t *testing.T) {
	episode := 1

	client := &fakeTraceMoeClient{
		result: &tracemoe.SearchResponse{
			FrameCount: 1,
			Result: []tracemoe.SearchResult{
				{
					AnilistID:  20458,
					Filename:   "test-anime.mp4",
					Episode:    &episode,
					From:       1122.16,
					To:         1126.12,
					Similarity: 0.9693,
					Video:      "https://example.com/video",
					Image:      "https://example.com/image",
				},
			},
		},
	}

	service := NewAnimeService(
		&fakeAniListClient{},
		client,
	)

	result, err := service.IdentifyAnime(
		context.Background(),
		"https://example.com/screenshot.jpg",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(result.Data),
		)
	}

	match := result.Data[0]

	if match.AnilistID != 20458 {
		t.Errorf(
			"expected anilist_id 20458, got %d",
			match.AnilistID,
		)
	}

	if match.Episode == nil || *match.Episode != 1 {
		t.Errorf("expected episode 1")
	}

	if match.Similarity != 0.9693 {
		t.Errorf(
			"expected similarity 0.9693, got %f",
			match.Similarity,
		)
	}
}

func TestAnimeService_IdentifyAnime_Error(t *testing.T) {
	expectedErr := errors.New("trace.moe unavailable")

	client := &fakeTraceMoeClient{
		err: expectedErr,
	}

	service := NewAnimeService(
		&fakeAniListClient{},
		client,
	)

	_, err := service.IdentifyAnime(
		context.Background(),
		"https://example.com/screenshot.jpg",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func intPtr(value int) *int {
	return &value
}
