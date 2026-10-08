package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type LibraryHandler struct {
	service service.LibraryService
}

func NewLibraryHandler(
	service service.LibraryService,
) *LibraryHandler {
	return &LibraryHandler{
		service: service,
	}
}

type AddLibraryAnimeRequest struct {
	ExternalID string `json:"external_id" binding:"required"`
}

type UpdateLibraryAnimeRequest struct {
	Status         *string  `json:"status"`
	CurrentEpisode *int     `json:"current_episode"`
	Rating         *float64 `json:"rating"`
	Notes          *string  `json:"notes"`
	StartedAt      *string  `json:"started_at"`
	CompletedAt    *string  `json:"completed_at"`
}

func (h *LibraryHandler) AddAnime(c *gin.Context) {
	var request AddLibraryAnimeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "external_id is required",
			},
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	result, err := h.service.AddAnime(
		c.Request.Context(),
		userIDInt64,
		request.ExternalID,
	)

	if err != nil {
		if errors.Is(err, service.ErrLibraryAnimeExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "ANIME_ALREADY_IN_LIBRARY",
					"message": "anime already exists in library",
				},
			})
			return
		}

		if errors.Is(err, service.ErrAnimeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to add anime to library",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"library_anime": result,
	})
}

func (h *LibraryHandler) GetLibrary(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	result, err := h.service.GetLibrary(
		c.Request.Context(),
		userIDInt64,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to fetch library",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *LibraryHandler) GetAnime(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid library anime id",
			},
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	result, err := h.service.GetAnime(
		c.Request.Context(),
		userIDInt64,
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrLibraryAnimeNotFound,
		) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found in library",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to fetch library anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"library_anime": result,
	})
}

func (h *LibraryHandler) UpdateAnime(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid library anime id",
			},
		})
		return
	}

	var request UpdateLibraryAnimeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	if request.Status != nil &&
		!isValidLibraryStatus(*request.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid status",
			},
		})
		return
	}

	if request.CurrentEpisode != nil &&
		*request.CurrentEpisode < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "current_episode cannot be negative",
			},
		})
		return
	}

	if request.Rating != nil &&
		(*request.Rating < 1 || *request.Rating > 10) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "rating must be between 1 and 10",
			},
		})
		return
	}

	if request.StartedAt != nil {
		if _, err := time.Parse(
			time.RFC3339,
			*request.StartedAt,
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "started_at must be a valid RFC3339 timestamp",
				},
			})
			return
		}
	}

	if request.CompletedAt != nil {
		if _, err := time.Parse(
			time.RFC3339,
			*request.CompletedAt,
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "completed_at must be a valid RFC3339 timestamp",
				},
			})
			return
		}
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	result, err := h.service.UpdateAnime(
		c.Request.Context(),
		userIDInt64,
		id,
		request.Status,
		request.CurrentEpisode,
		request.Rating,
		request.Notes,
		request.StartedAt,
		request.CompletedAt,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrLibraryAnimeNotFound,
		) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found in library",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to update library anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"library_anime": result,
	})
}

func isValidLibraryStatus(status string) bool {
	switch status {
	case "WATCHING",
		"COMPLETED",
		"PLAN_TO_WATCH",
		"ON_HOLD",
		"DROPPED":
		return true
	default:
		return false
	}
}

func (h *LibraryHandler) DeleteAnime(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid library anime id",
			},
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	userIDInt64, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "invalid user context",
			},
		})
		return
	}

	err = h.service.DeleteAnime(
		c.Request.Context(),
		userIDInt64,
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrLibraryAnimeNotFound,
		) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found in library",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to delete anime from library",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}
