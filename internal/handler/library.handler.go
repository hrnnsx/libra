package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		if err == service.ErrLibraryAnimeExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "ANIME_ALREADY_IN_LIBRARY",
					"message": "anime already exists in library",
				},
			})
			return
		}

		if err == service.ErrAnimeNotFound {
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
