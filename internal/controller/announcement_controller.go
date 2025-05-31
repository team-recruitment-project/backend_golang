package controller

import (
	"backend_golang/internal/controller/request"
	"backend_golang/internal/models"
	"backend_golang/internal/service"
	servicemodels "backend_golang/internal/service/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AnnouncementController interface {
	Announce(c *gin.Context)
	GetAnnouncement(c *gin.Context)
	GetAnnouncements(c *gin.Context)
}

type announcementController struct {
	announcementService service.AnnouncementService
}

func NewAnnouncementController(announcementService service.AnnouncementService) AnnouncementController {
	return &announcementController{
		announcementService: announcementService,
	}
}

func (a *announcementController) Announce(c *gin.Context) {
	memberID := c.Value("userID").(string)
	if memberID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req := &request.PostAnnouncement{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		validationErrors := make([]models.ValidationError, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.NewValidationError(err))
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	announcementID, err := a.announcementService.Announce(c, servicemodels.RegisterAnnouncement{
		TeamID:   req.TeamID,
		MemberID: memberID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"announcementID": announcementID})
}

func (a *announcementController) GetAnnouncement(c *gin.Context) {
	announcementID, err := strconv.Atoi(c.Param("announcementID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	announcement, err := a.announcementService.GetAnnouncement(c, announcementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, announcement)
}

func (a *announcementController) GetAnnouncements(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	log.Println("page", c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	size, err := strconv.Atoi(c.Query("size"))
	log.Println("size", c.Query("size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skillParams := c.Query("skill")
	var skills []string
	if skillParams != "" {
		skills = strings.Split(skillParams, ",")
	}

	positionParams := c.Query("position")
	var positions []string
	if positionParams != "" {
		positions = strings.Split(positionParams, ",")
	}

	keyword := c.Query("keyword")

	announcements, err := a.announcementService.GetAnnouncements(c, page, size, skills, positions, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, announcements)
}
