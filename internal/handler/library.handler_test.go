package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type fakeLibraryService struct {
	addResult *model.LibraryAnime
	addErr    error

	getLibraryResult []model.LibraryAnime
	getLibraryErr    error

	getAnimeResult *model.LibraryAnime
	getAnimeErr    error

	updateResult *model.LibraryAnime
	updateErr    error

	deleteErr error

	lastUserID     int64
	lastExternalID string
	lastAnimeID    int64

	lastStatus         *string
	lastCurrentEpisode *int
	lastRating         *float64
	lastNotes          *string
	lastStartedAt      *string
	lastCompletedAt    *string
}

func (f *fakeLibraryService) AddAnime(
	ctx context.Context,
	userID int64,
	externalID string,
) (*model.LibraryAnime, error) {
	f.lastUserID = userID
	f.lastExternalID = externalID

	return f.addResult, f.addErr
}

func (f *fakeLibraryService) GetLibrary(
	ctx context.Context,
	userID int64,
) ([]model.LibraryAnime, error) {
	f.lastUserID = userID

	return f.getLibraryResult, f.getLibraryErr
}

func (f *fakeLibraryService) GetAnime(
	ctx context.Context,
	userID int64,
	id int64,
) (*model.LibraryAnime, error) {
	f.lastUserID = userID
	f.lastAnimeID = id

	return f.getAnimeResult, f.getAnimeErr
}

func (f *fakeLibraryService) UpdateAnime(
	ctx context.Context,
	userID int64,
	id int64,
	status *string,
	currentEpisode *int,
	rating *float64,
	notes *string,
	startedAt *string,
	completedAt *string,
) (*model.LibraryAnime, error) {
	f.lastUserID = userID
	f.lastAnimeID = id
	f.lastStatus = status
	f.lastCurrentEpisode = currentEpisode
	f.lastRating = rating
	f.lastNotes = notes
	f.lastStartedAt = startedAt
	f.lastCompletedAt = completedAt

	return f.updateResult, f.updateErr
}

func (f *fakeLibraryService) DeleteAnime(
	ctx context.Context,
	userID int64,
	id int64,
) error {
	f.lastUserID = userID
	f.lastAnimeID = id

	return f.deleteErr
}

func setupLibraryHandler(
	libraryService service.LibraryService,
) (*gin.Engine, *LibraryHandler) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewLibraryHandler(libraryService)

	return router, handler
}

func setUserID(c *gin.Context, userID int64) {
	c.Set("user_id", userID)
}

func TestLibraryHandler_AddAnime(t *testing.T) {
	fakeService := &fakeLibraryService{
		addResult: &model.LibraryAnime{
			ID:             1,
			UserID:         10,
			AnimeID:        100,
			Status:         "PLAN_TO_WATCH",
			CurrentEpisode: 0,
		},
	}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.AddAnime(c)
		},
	)

	body := `{
		"external_id": "20464"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			recorder.Code,
		)
	}

	var response struct {
		LibraryAnime model.LibraryAnime `json:"library_anime"`
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

	if response.LibraryAnime.ID != 1 {
		t.Errorf(
			"expected id 1, got %d",
			response.LibraryAnime.ID,
		)
	}

	if fakeService.lastUserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			fakeService.lastUserID,
		)
	}

	if fakeService.lastExternalID != "20464" {
		t.Errorf(
			"expected external_id 20464, got %s",
			fakeService.lastExternalID,
		)
	}
}

func TestLibraryHandler_AddAnime_InvalidRequest(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.AddAnime(c)
		},
	)

	body := `{}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_AddAnime_AlreadyExists(t *testing.T) {
	fakeService := &fakeLibraryService{
		addErr: service.ErrLibraryAnimeExists,
	}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.AddAnime(c)
		},
	)

	body := `{
		"external_id": "20464"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusConflict {
		t.Fatalf(
			"expected status 409, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_AddAnime_NotFound(t *testing.T) {
	fakeService := &fakeLibraryService{
		addErr: service.ErrAnimeNotFound,
	}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.AddAnime(c)
		},
	)

	body := `{
		"external_id": "999999"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_AddAnime_ServiceError(t *testing.T) {
	fakeService := &fakeLibraryService{
		addErr: errors.New("database unavailable"),
	}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.AddAnime(c)
		},
	)

	body := `{
		"external_id": "20464"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_AddAnime_MissingUserContext(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.POST(
		"/library/animes",
		handler.AddAnime,
	)

	body := `{
		"external_id": "20464"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/library/animes",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_GetLibrary(t *testing.T) {
	fakeService := &fakeLibraryService{
		getLibraryResult: []model.LibraryAnime{
			{
				ID:             1,
				UserID:         10,
				AnimeID:        100,
				Status:         "WATCHING",
				CurrentEpisode: 5,
			},
			{
				ID:             2,
				UserID:         10,
				AnimeID:        200,
				Status:         "COMPLETED",
				CurrentEpisode: 12,
			},
		},
	}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetLibrary(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes",
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
		Data []model.LibraryAnime `json:"data"`
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

	if len(response.Data) != 2 {
		t.Fatalf(
			"expected 2 anime, got %d",
			len(response.Data),
		)
	}

	if fakeService.lastUserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			fakeService.lastUserID,
		)
	}
}

func TestLibraryHandler_GetLibrary_ServiceError(t *testing.T) {
	fakeService := &fakeLibraryService{
		getLibraryErr: errors.New("database unavailable"),
	}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetLibrary(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_GetAnime(t *testing.T) {
	fakeService := &fakeLibraryService{
		getAnimeResult: &model.LibraryAnime{
			ID:             1,
			UserID:         10,
			AnimeID:        100,
			Status:         "WATCHING",
			CurrentEpisode: 5,
		},
	}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes/1",
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
		LibraryAnime model.LibraryAnime `json:"library_anime"`
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

	if response.LibraryAnime.ID != 1 {
		t.Errorf(
			"expected id 1, got %d",
			response.LibraryAnime.ID,
		)
	}

	if fakeService.lastAnimeID != 1 {
		t.Errorf(
			"expected anime id 1, got %d",
			fakeService.lastAnimeID,
		)
	}
}

func TestLibraryHandler_GetAnime_InvalidID(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes/abc",
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

func TestLibraryHandler_GetAnime_NotFound(t *testing.T) {
	fakeService := &fakeLibraryService{
		getAnimeErr: repository.ErrLibraryAnimeNotFound,
	}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes/999",
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

func TestLibraryHandler_GetAnime_ServiceError(t *testing.T) {
	fakeService := &fakeLibraryService{
		getAnimeErr: errors.New("database unavailable"),
	}

	router, handler := setupLibraryHandler(fakeService)

	router.GET(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.GetAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/library/animes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime(t *testing.T) {
	fakeService := &fakeLibraryService{
		updateResult: &model.LibraryAnime{
			ID:             1,
			UserID:         10,
			AnimeID:        100,
			Status:         "COMPLETED",
			CurrentEpisode: 12,
		},
	}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"status": "COMPLETED",
		"current_episode": 12,
		"rating": 9.5,
		"notes": "Great anime",
		"started_at": "2026-10-01T10:00:00Z",
		"completed_at": "2026-10-05T10:00:00Z"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/1",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	if fakeService.lastUserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			fakeService.lastUserID,
		)
	}

	if fakeService.lastAnimeID != 1 {
		t.Errorf(
			"expected anime id 1, got %d",
			fakeService.lastAnimeID,
		)
	}

	if fakeService.lastStatus == nil ||
		*fakeService.lastStatus != "COMPLETED" {
		t.Errorf("expected status COMPLETED")
	}

	if fakeService.lastCurrentEpisode == nil ||
		*fakeService.lastCurrentEpisode != 12 {
		t.Errorf("expected current_episode 12")
	}

	if fakeService.lastRating == nil ||
		*fakeService.lastRating != 9.5 {
		t.Errorf("expected rating 9.5")
	}

	if fakeService.lastNotes == nil ||
		*fakeService.lastNotes != "Great anime" {
		t.Errorf("expected notes")
	}
}

func TestLibraryHandler_UpdateAnime_InvalidID(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"status": "COMPLETED"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/abc",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime_InvalidStatus(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"status": "INVALID"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/1",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime_NegativeEpisode(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"current_episode": -1
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/1",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime_InvalidRating(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"rating": 11
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/1",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime_InvalidStartedAt(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"started_at": "not-a-date"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/1",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_UpdateAnime_NotFound(t *testing.T) {
	fakeService := &fakeLibraryService{
		updateErr: repository.ErrLibraryAnimeNotFound,
	}

	router, handler := setupLibraryHandler(fakeService)

	router.PATCH(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.UpdateAnime(c)
		},
	)

	body := `{
		"status": "COMPLETED"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/library/animes/999",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			recorder.Code,
		)
	}
}

func TestLibraryHandler_DeleteAnime(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.DELETE(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.DeleteAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/library/animes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status 204, got %d",
			recorder.Code,
		)
	}

	if fakeService.lastUserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			fakeService.lastUserID,
		)
	}

	if fakeService.lastAnimeID != 1 {
		t.Errorf(
			"expected anime id 1, got %d",
			fakeService.lastAnimeID,
		)
	}

	if recorder.Body.Len() != 0 {
		t.Errorf(
			"expected empty response body, got %q",
			recorder.Body.String(),
		)
	}
}

func TestLibraryHandler_DeleteAnime_InvalidID(t *testing.T) {
	fakeService := &fakeLibraryService{}

	router, handler := setupLibraryHandler(fakeService)

	router.DELETE(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.DeleteAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/library/animes/abc",
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

func TestLibraryHandler_DeleteAnime_NotFound(t *testing.T) {
	fakeService := &fakeLibraryService{
		deleteErr: repository.ErrLibraryAnimeNotFound,
	}

	router, handler := setupLibraryHandler(fakeService)

	router.DELETE(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.DeleteAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/library/animes/999",
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

func TestLibraryHandler_DeleteAnime_ServiceError(t *testing.T) {
	fakeService := &fakeLibraryService{
		deleteErr: errors.New("database unavailable"),
	}

	router, handler := setupLibraryHandler(fakeService)

	router.DELETE(
		"/library/animes/:id",
		func(c *gin.Context) {
			setUserID(c, 10)
			handler.DeleteAnime(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/library/animes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			recorder.Code,
		)
	}
}
