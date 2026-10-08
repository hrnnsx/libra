package handler

import (
	"context"
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

type fakeGroupLibraryAnimeService struct {
	addAnimeFunc func(
		ctx context.Context,
		userID int64,
		groupID int64,
		libraryAnimeID int64,
	) (*model.GroupLibraryAnime, error)

	getAnimesFunc func(
		ctx context.Context,
		userID int64,
		groupID int64,
	) ([]model.GroupLibraryAnime, error)

	deleteAnimeFunc func(
		ctx context.Context,
		userID int64,
		groupID int64,
		libraryAnimeID int64,
	) error
}

func (f *fakeGroupLibraryAnimeService) AddAnime(
	ctx context.Context,
	userID int64,
	groupID int64,
	libraryAnimeID int64,
) (*model.GroupLibraryAnime, error) {
	if f.addAnimeFunc != nil {
		return f.addAnimeFunc(
			ctx,
			userID,
			groupID,
			libraryAnimeID,
		)
	}

	return nil, nil
}

func (f *fakeGroupLibraryAnimeService) GetAnimes(
	ctx context.Context,
	userID int64,
	groupID int64,
) ([]model.GroupLibraryAnime, error) {
	if f.getAnimesFunc != nil {
		return f.getAnimesFunc(
			ctx,
			userID,
			groupID,
		)
	}

	return nil, nil
}

func (f *fakeGroupLibraryAnimeService) DeleteAnime(
	ctx context.Context,
	userID int64,
	groupID int64,
	libraryAnimeID int64,
) error {
	if f.deleteAnimeFunc != nil {
		return f.deleteAnimeFunc(
			ctx,
			userID,
			groupID,
			libraryAnimeID,
		)
	}

	return nil
}

func setGroupLibraryAnimeUserID(c *gin.Context, userID int64) {
	c.Set("user_id", userID)
}

func TestGroupLibraryAnimeHandler_AddAnime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			addAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) (*model.GroupLibraryAnime, error) {
				return &model.GroupLibraryAnime{
					GroupID:        groupID,
					LibraryAnimeID: libraryAnimeID,
				}, nil
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				rec.Code,
			)
		}

		if !strings.Contains(
			rec.Body.String(),
			`"group_id":1`,
		) {
			t.Errorf("expected group_id in response")
		}

		if !strings.Contains(
			rec.Body.String(),
			`"library_anime_id":100`,
		) {
			t.Errorf("expected library_anime_id in response")
		}
	})

	t.Run("invalid group id", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", handler.AddAnime)

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/abc/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("invalid library anime id", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", handler.AddAnime)

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":0}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", handler.AddAnime)

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("missing user context", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", handler.AddAnime)

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("invalid user context type", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			c.Set("user_id", "10")
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			addAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) (*model.GroupLibraryAnime, error) {
				return nil, service.ErrGroupNotFoundForUser
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})

	t.Run("library anime not found", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			addAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) (*model.GroupLibraryAnime, error) {
				return nil, service.ErrLibraryAnimeNotFound
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})

	t.Run("anime already exists", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			addAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) (*model.GroupLibraryAnime, error) {
				return nil, service.ErrGroupLibraryAnimeExists
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusConflict,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("database error")

		service := &fakeGroupLibraryAnimeService{
			addAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) (*model.GroupLibraryAnime, error) {
				return nil, expectedErr
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.POST("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.AddAnime(c)
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/groups/1/animes",
			strings.NewReader(`{"library_anime_id":100}`),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				rec.Code,
			)
		}
	})
}

func TestGroupLibraryAnimeHandler_GetAnimes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			getAnimesFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
			) ([]model.GroupLibraryAnime, error) {
				return []model.GroupLibraryAnime{
					{
						GroupID:        1,
						LibraryAnimeID: 100,
					},
					{
						GroupID:        1,
						LibraryAnimeID: 101,
					},
				}, nil
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.GetAnimes(c)
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/1/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

		if !strings.Contains(
			rec.Body.String(),
			`"library_anime_id":100`,
		) {
			t.Errorf("expected library anime 100")
		}

		if !strings.Contains(
			rec.Body.String(),
			`"library_anime_id":101`,
		) {
			t.Errorf("expected library anime 101")
		}
	})

	t.Run("invalid group id", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", handler.GetAnimes)

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/abc/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("missing user context", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", handler.GetAnimes)

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/1/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("invalid user context type", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", func(c *gin.Context) {
			c.Set("user_id", "10")
			handler.GetAnimes(c)
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/1/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			getAnimesFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
			) ([]model.GroupLibraryAnime, error) {
				return nil, service.ErrGroupNotFoundForUser
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.GetAnimes(c)
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/1/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("database error")

		service := &fakeGroupLibraryAnimeService{
			getAnimesFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
			) ([]model.GroupLibraryAnime, error) {
				return nil, expectedErr
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.GET("/groups/:id/animes", func(c *gin.Context) {
			setGroupLibraryAnimeUserID(c, 10)
			handler.GetAnimes(c)
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/groups/1/animes",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				rec.Code,
			)
		}
	})
}

func TestGroupLibraryAnimeHandler_DeleteAnime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		var receivedUserID int64
		var receivedGroupID int64
		var receivedLibraryAnimeID int64

		service := &fakeGroupLibraryAnimeService{
			deleteAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) error {
				receivedUserID = userID
				receivedGroupID = groupID
				receivedLibraryAnimeID = libraryAnimeID
				return nil
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			func(c *gin.Context) {
				setGroupLibraryAnimeUserID(c, 10)
				handler.DeleteAnime(c)
			},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNoContent,
				rec.Code,
			)
		}

		if receivedUserID != 10 {
			t.Errorf(
				"expected user ID 10, got %d",
				receivedUserID,
			)
		}

		if receivedGroupID != 1 {
			t.Errorf(
				"expected group ID 1, got %d",
				receivedGroupID,
			)
		}

		if receivedLibraryAnimeID != 100 {
			t.Errorf(
				"expected library anime ID 100, got %d",
				receivedLibraryAnimeID,
			)
		}
	})

	t.Run("invalid group id", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			handler.DeleteAnime,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/abc/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("invalid library anime id", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			handler.DeleteAnime,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/abc",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("missing user context", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			handler.DeleteAnime,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("invalid user context type", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			func(c *gin.Context) {
				c.Set("user_id", "10")
				handler.DeleteAnime(c)
			},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			deleteAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) error {
				return service.ErrGroupNotFoundForUser
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			func(c *gin.Context) {
				setGroupLibraryAnimeUserID(c, 10)
				handler.DeleteAnime(c)
			},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})

	t.Run("anime not found in group", func(t *testing.T) {
		service := &fakeGroupLibraryAnimeService{
			deleteAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) error {
				return repository.ErrGroupLibraryAnimeNotFound
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			func(c *gin.Context) {
				setGroupLibraryAnimeUserID(c, 10)
				handler.DeleteAnime(c)
			},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("database error")

		service := &fakeGroupLibraryAnimeService{
			deleteAnimeFunc: func(
				ctx context.Context,
				userID int64,
				groupID int64,
				libraryAnimeID int64,
			) error {
				return expectedErr
			},
		}

		handler := NewGroupLibraryAnimeHandler(service)

		router := gin.New()
		router.DELETE(
			"/groups/:id/animes/:library_anime_id",
			func(c *gin.Context) {
				setGroupLibraryAnimeUserID(c, 10)
				handler.DeleteAnime(c)
			},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/groups/1/animes/100",
			nil,
		)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				rec.Code,
			)
		}
	})
}
