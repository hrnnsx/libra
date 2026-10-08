package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type GroupHandler struct {
	service service.GroupService
}

func NewGroupHandler(
	service service.GroupService,
) *GroupHandler {
	return &GroupHandler{
		service: service,
	}
}

type CreateGroupRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description *string `json:"description"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var request CreateGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
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

	result, err := h.service.CreateGroup(
		c.Request.Context(),
		userIDInt64,
		request.Name,
		request.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to create group",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"group": result,
	})
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
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

	result, err := h.service.GetGroups(
		c.Request.Context(),
		userIDInt64,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "failed to fetch groups",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
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

	result, err := h.service.GetGroup(
		c.Request.Context(),
		userIDInt64,
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrGroupNotFound,
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
				"message": "failed to fetch group",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group": result,
	})
}

func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid group id",
			},
		})
		return
	}

	var request UpdateGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid request body",
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

	result, err := h.service.UpdateGroup(
		c.Request.Context(),
		userIDInt64,
		id,
		request.Name,
		request.Description,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrGroupNotFound,
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
				"message": "failed to update group",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group": result,
	})
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id < 1 {
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

	err = h.service.DeleteGroup(
		c.Request.Context(),
		userIDInt64,
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrGroupNotFound,
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
				"message": "failed to delete group",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}
