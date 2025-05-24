package controller

import (
	"backend_golang/internal/controller/request"
	"backend_golang/internal/models"
	"backend_golang/internal/service"
	servicemodels "backend_golang/internal/service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AnnouncementController interface {
	Announce(c *gin.Context)
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
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"announcementID": announcementID})
}
