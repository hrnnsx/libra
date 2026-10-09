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

// AddAnime godoc
// @Summary Add anime to group
// @Description Add an anime from the user's library to a group.
// @Tags Group Anime
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param request body AddGroupAnimeRequest true "Library anime to add"
// @Success 201 {object} map[string]interface{} "Anime added to group"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 404 {object} map[string]interface{} "Group or library anime not found"
// @Failure 409 {object} map[string]interface{} "Anime already exists in group"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /me/groups/{id}/animes [post]
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

// GetAnimes godoc
// @Summary Get anime in group
// @Description Get all library anime assigned to a specific group.
// @Tags Group Anime
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]interface{} "Group anime retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid group ID"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /me/groups/{id}/animes [get]
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

// DeleteAnime godoc
// @Summary Remove anime from group
// @Description Remove an anime from a group without deleting it from the user's library.
// @Tags Group Anime
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param library_anime_id path int true "Library anime ID"
// @Success 204 "Anime removed from group"
// @Failure 400 {object} map[string]interface{} "Invalid group ID or library anime ID"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 404 {object} map[string]interface{} "Group or anime in group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /me/groups/{id}/animes/{library_anime_id} [delete]
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
