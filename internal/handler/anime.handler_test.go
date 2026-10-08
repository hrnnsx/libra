package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/service"
)

type fakeAnimeService struct {
	browseResult *service.AnimeListResponse
	browseErr    error

	searchResult *service.AnimeListResponse
	searchErr    error

	getResult *service.AnimeResponse
	getErr    error

	identifyResult *service.AnimeIdentificationListResponse
	identifyErr    error
}

func (f *fakeAnimeService) BrowseAnime(
	ctx context.Context,
	page int,
	perPage int,
) (*service.AnimeListResponse, error) {
	return f.browseResult, f.browseErr
}

func (f *fakeAnimeService) SearchAnime(
	ctx context.Context,
	params service.AnimeSearchParams,
) (*service.AnimeListResponse, error) {
	return f.searchResult, f.searchErr
}

func (f *fakeAnimeService) GetAnime(
	ctx context.Context,
	externalID string,
) (*service.AnimeResponse, error) {
	return f.getResult, f.getErr
}

func (f *fakeAnimeService) IdentifyAnime(
	ctx context.Context,
	imageURL string,
) (*service.AnimeIdentificationListResponse, error) {
	return f.identifyResult, f.identifyErr
}

func setupAnimeHandler(
	service service.AnimeService,
) (*gin.Engine, *AnimeHandler) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewAnimeHandler(service)

	return router, handler
}

func TestAnimeHandler_BrowseAnime(t *testing.T) {
	fakeService := &fakeAnimeService{
		browseResult: &service.AnimeListResponse{
			Data: []service.AnimeResponse{
				{
					ExternalID: "20464",
					Title: service.AnimeTitle{
						Romaji:  "Mahouka Koukou no Rettousei",
						English: "The Irregular at Magic High School",
					},
				},
			},
			Pagination: service.Pagination{
				Page:        1,
				PerPage:     20,
				HasNextPage: true,
				Total:       100,
			},
		},
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes", handler.BrowseAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var response service.AnimeListResponse

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(response.Data) != 1 {
		t.Fatalf(
			"expected 1 anime, got %d",
			len(response.Data),
		)
	}

	if response.Data[0].ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			response.Data[0].ExternalID,
		)
	}
}

func TestAnimeHandler_BrowseAnime_Pagination(t *testing.T) {
	fakeService := &fakeAnimeService{
		browseResult: &service.AnimeListResponse{
			Data: []service.AnimeResponse{},
			Pagination: service.Pagination{
				Page:    2,
				PerPage: 5,
			},
		},
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes", handler.BrowseAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes?page=2&per_page=5",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_BrowseAnime_InvalidPagination(t *testing.T) {
	fakeService := &fakeAnimeService{}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes", handler.BrowseAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes?page=abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_BrowseAnime_ExternalAPIError(t *testing.T) {
	fakeService := &fakeAnimeService{
		browseErr: errors.New("anilist unavailable"),
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes", handler.BrowseAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status 502, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_SearchAnime(t *testing.T) {
	fakeService := &fakeAnimeService{
		searchResult: &service.AnimeListResponse{
			Data: []service.AnimeResponse{
				{
					ExternalID: "20464",
					Title: service.AnimeTitle{
						Romaji: "Mahouka Koukou no Rettousei",
					},
				},
			},
			Pagination: service.Pagination{
				Page:        1,
				PerPage:     20,
				HasNextPage: false,
				Total:       1,
			},
		},
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/search", handler.SearchAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/search?title=mahouka",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var response service.AnimeListResponse

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(response.Data) != 1 {
		t.Fatalf(
			"expected 1 anime, got %d",
			len(response.Data),
		)
	}

	if response.Data[0].ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			response.Data[0].ExternalID,
		)
	}
}

func TestAnimeHandler_SearchAnime_InvalidSeasonYear(t *testing.T) {
	fakeService := &fakeAnimeService{}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/search", handler.SearchAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/search?season_year=abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_SearchAnime_ServiceError(t *testing.T) {
	fakeService := &fakeAnimeService{
		searchErr: errors.New("anilist unavailable"),
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/search", handler.SearchAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/search?title=mahouka",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status 502, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_GetAnime(t *testing.T) {
	fakeService := &fakeAnimeService{
		getResult: &service.AnimeResponse{
			ExternalID: "20464",
			Title: service.AnimeTitle{
				Romaji:  "Mahouka Koukou no Rettousei",
				English: "The Irregular at Magic High School",
			},
		},
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/:external_id", handler.GetAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/20464",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var response struct {
		Anime service.AnimeResponse `json:"anime"`
	}

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response.Anime.ExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			response.Anime.ExternalID,
		)
	}
}

func TestAnimeHandler_GetAnime_InvalidExternalID(t *testing.T) {
	fakeService := &fakeAnimeService{}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/:external_id", handler.GetAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_GetAnime_NotFound(t *testing.T) {
	fakeService := &fakeAnimeService{
		getErr: service.ErrAnimeNotFound,
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/:external_id", handler.GetAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/999999",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_GetAnime_ServiceError(t *testing.T) {
	fakeService := &fakeAnimeService{
		getErr: errors.New("anilist unavailable"),
	}

	router, handler := setupAnimeHandler(fakeService)
	router.GET("/animes/:external_id", handler.GetAnime)

	req := httptest.NewRequest(
		http.MethodGet,
		"/animes/20464",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status 502, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_IdentifyAnime(t *testing.T) {
	fakeService := &fakeAnimeService{
		identifyResult: &service.AnimeIdentificationListResponse{
			Data: []service.AnimeIdentificationResponse{
				{
					AnilistID:  20458,
					Filename:   "test-anime.mp4",
					Episode:    intPtr(1),
					From:       1122.16,
					To:         1126.12,
					Similarity: 0.9693,
					Video:      "https://example.com/video",
					Image:      "https://example.com/image",
				},
			},
		},
	}

	router, handler := setupAnimeHandler(fakeService)
	router.POST("/animes/identify", handler.IdentifyAnime)

	body := `{
		"url": "https://example.com/screenshot.jpg"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/animes/identify",
		stringReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	var response service.AnimeIdentificationListResponse

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if len(response.Data) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(response.Data),
		)
	}

	if response.Data[0].AnilistID != 20458 {
		t.Errorf(
			"expected anilist_id 20458, got %d",
			response.Data[0].AnilistID,
		)
	}

	if response.Data[0].Similarity != 0.9693 {
		t.Errorf(
			"expected similarity 0.9693, got %f",
			response.Data[0].Similarity,
		)
	}
}

func TestAnimeHandler_IdentifyAnime_InvalidURL(t *testing.T) {
	fakeService := &fakeAnimeService{}

	router, handler := setupAnimeHandler(fakeService)
	router.POST("/animes/identify", handler.IdentifyAnime)

	body := `{
		"url": ""
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/animes/identify",
		stringReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_IdentifyAnime_InvalidBody(t *testing.T) {
	fakeService := &fakeAnimeService{}

	router, handler := setupAnimeHandler(fakeService)
	router.POST("/animes/identify", handler.IdentifyAnime)

	body := `{
		"url":
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/animes/identify",
		stringReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestAnimeHandler_IdentifyAnime_ServiceError(t *testing.T) {
	fakeService := &fakeAnimeService{
		identifyErr: errors.New("trace.moe unavailable"),
	}

	router, handler := setupAnimeHandler(fakeService)
	router.POST("/animes/identify", handler.IdentifyAnime)

	body := `{
		"url": "https://example.com/screenshot.jpg"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/animes/identify",
		stringReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status 502, got %d",
			recorder.Code,
		)
	}
}

func stringReader(value string) *stringReaderType {
	return &stringReaderType{
		value: value,
	}
}

type stringReaderType struct {
	value string
	index int
}

func (r *stringReaderType) Read(p []byte) (int, error) {
	if r.index >= len(r.value) {
		return 0, errors.New("EOF")
	}

	n := copy(p, r.value[r.index:])
	r.index += n

	return n, nil
}

func intPtr(value int) *int {
	return &value
}
