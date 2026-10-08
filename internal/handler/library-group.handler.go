package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type GroupLibraryAnimeHandler struct {
	service service.GroupLibraryAnimeService
}

func NewGroupLibraryAnimeHandler(
	service service.GroupLibraryAnimeService,
) *GroupLibraryAnimeHandler {
	return &GroupLibraryAnimeHandler{
		service: service,
	}
}

type AddGroupAnimeRequest struct {
	LibraryAnimeID int64 `json:"library_anime_id" binding:"required,min=1"`
}

func (h *GroupLibraryAnimeHandler) AddAnime(c *gin.Context) {
	groupID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || groupID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid group id",
			},
		})
		return
	}

	var request AddGroupAnimeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "library_anime_id is required",
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
		groupID,
		request.LibraryAnimeID,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrGroupNotFoundForUser,
		):
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "GROUP_NOT_FOUND",
					"message": "group not found",
				},
			})

		case errors.Is(
			err,
			service.ErrLibraryAnimeNotFound,
		):
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found in library",
				},
			})

		case errors.Is(
			err,
			service.ErrGroupLibraryAnimeExists,
		):
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "ANIME_ALREADY_IN_GROUP",
					"message": "anime already exists in group",
				},
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "failed to add anime to group",
				},
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"group_anime": result,
	})
}

func (h *GroupLibraryAnimeHandler) GetAnimes(c *gin.Context) {
	groupID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || groupID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid group id",
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

	result, err := h.service.GetAnimes(
		c.Request.Context(),
		userIDInt64,
		groupID,
	)

	if err != nil {
		if errors.Is(
			err,
			service.ErrGroupNotFoundForUser,
		) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "GROUP_NOT_FOUND",
					"message": "group not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to fetch group anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *GroupLibraryAnimeHandler) DeleteAnime(c *gin.Context) {
	groupID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || groupID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid group id",
			},
		})
		return
	}

	libraryAnimeID, err := strconv.ParseInt(
		c.Param("library_anime_id"),
		10,
		64,
	)

	if err != nil || libraryAnimeID < 1 {
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
		groupID,
		libraryAnimeID,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			service.ErrGroupNotFoundForUser,
		):
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "GROUP_NOT_FOUND",
					"message": "group not found",
				},
			})

		case errors.Is(
			err,
			repository.ErrGroupLibraryAnimeNotFound,
		):
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found in group",
				},
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "failed to remove anime from group",
				},
			})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
