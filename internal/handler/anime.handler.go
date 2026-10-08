package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/service"
)

type AnimeHandler struct {
	service service.AnimeService
}

func NewAnimeHandler(service service.AnimeService) *AnimeHandler {
	return &AnimeHandler{
		service: service,
	}
}

func parsePagination(c *gin.Context) (int, int, error) {
	page := 1
	perPage := 20

	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return 0, 0, err
		}

		page = parsed
	}

	if value := c.Query("per_page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			return 0, 0, err
		}

		perPage = parsed
	}

	return page, perPage, nil
}

func (h *AnimeHandler) BrowseAnime(c *gin.Context) {
	page, perPage, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid pagination parameters",
			},
		})
		return
	}

	result, err := h.service.BrowseAnime(
		c.Request.Context(),
		page,
		perPage,
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "EXTERNAL_API_ERROR",
				"message": "failed to fetch anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AnimeHandler) SearchAnime(c *gin.Context) {
	page, perPage, err := parsePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "invalid pagination parameters",
			},
		})
		return
	}

	seasonYear := 0

	if value := c.Query("season_year"); value != "" {
		seasonYear, err = strconv.Atoi(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "season_year must be a number",
				},
			})
			return
		}
	}

	result, err := h.service.SearchAnime(
		c.Request.Context(),
		service.AnimeSearchParams{
			Title:      c.Query("title"),
			ExternalID: c.Query("external_id"),
			Season:     c.Query("season"),
			SeasonYear: seasonYear,
			Genre:      c.Query("genre"),
			Status:     c.Query("status"),
			Format:     c.Query("format"),
			Page:       page,
			PerPage:    perPage,
		},
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "EXTERNAL_API_ERROR",
				"message": "failed to search anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AnimeHandler) GetAnime(c *gin.Context) {
	externalID := c.Param("external_id")

	if _, err := strconv.Atoi(externalID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "external_id must be a number",
			},
		})
		return
	}

	result, err := h.service.GetAnime(
		c.Request.Context(),
		externalID,
	)
	if err != nil {
		if err == service.ErrAnimeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ANIME_NOT_FOUND",
					"message": "anime not found",
				},
			})
			return
		}

		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "EXTERNAL_API_ERROR",
				"message": "failed to fetch anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"anime": result,
	})
}
