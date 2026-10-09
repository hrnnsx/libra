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

type IdentifyAnimeRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// BrowseAnime godoc
// @Summary Browse anime
// @Description Get a paginated list of popular anime.
// @Tags Anime
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param per_page query int false "Items per page" default(20) minimum(1) maximum(50)
// @Success 200 {object} service.AnimeListResponse
// @Failure 400 {object} map[string]interface{} "Invalid pagination parameters"
// @Failure 502 {object} map[string]interface{} "External API error"
// @Router /animes [get]
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

// SearchAnime godoc
// @Summary Search anime
// @Description Search anime by title, external ID, season, year, genre, status, or format.
// @Tags Anime
// @Produce json
// @Param title query string false "Anime title"
// @Param external_id query string false "AniList anime ID"
// @Param season query string false "Season" Enums(WINTER,SPRING,SUMMER,FALL)
// @Param season_year query int false "Season year"
// @Param genre query string false "Anime genre"
// @Param status query string false "Anime status" Enums(FINISHED,AIRING,NOT_YET_RELEASED,CANCELLED)
// @Param format query string false "Anime format" Enums(TV,TV_SHORT,MOVIE,SPECIAL,OVA,ONA,MUSIC)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param per_page query int false "Items per page" default(20) minimum(1) maximum(50)
// @Success 200 {object} service.AnimeListResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 502 {object} map[string]interface{} "External API error"
// @Router /animes/search [get]
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

// IdentifyAnime godoc
// @Summary Identify anime from image
// @Description Identify an anime scene from an image URL using trace.moe.
// @Tags Anime
// @Accept json
// @Produce json
// @Param request body IdentifyAnimeRequest true "Image URL"
// @Success 200 {object} service.AnimeIdentificationListResponse
// @Failure 400 {object} map[string]interface{} "Invalid image URL"
// @Failure 502 {object} map[string]interface{} "External API error"
// @Router /animes/identify [post]
func (h *AnimeHandler) IdentifyAnime(c *gin.Context) {
	var request IdentifyAnimeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "valid image url is required",
			},
		})
		return
	}

	result, err := h.service.IdentifyAnime(
		c.Request.Context(),
		request.URL,
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "EXTERNAL_API_ERROR",
				"message": "failed to identify anime",
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAnime godoc
// @Summary Get anime detail
// @Description Get anime details by AniList external ID.
// @Tags Anime
// @Produce json
// @Param external_id path string true "AniList anime ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid external ID"
// @Failure 404 {object} map[string]interface{} "Anime not found"
// @Failure 502 {object} map[string]interface{} "External API error"
// @Router /animes/{external_id} [get]
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
